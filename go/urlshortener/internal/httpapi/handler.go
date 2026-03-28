package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"urlshortener/internal/domain"
)

type shortener interface {
	Shorten(ctx context.Context, longURL domain.LongURL) (domain.ShortURL, error)
}

type resolver interface {
	Resolve(ctx context.Context, shortURL domain.ShortURL) (domain.LongURL, error)
}

type Handler struct {
	shortener shortener
	resolver  resolver
	baseURL   string
}

func NewHandler(shortener shortener, resolver resolver, baseURL string) *Handler {
	return &Handler{
		shortener: shortener,
		resolver:  resolver,
		baseURL:   strings.TrimRight(baseURL, "/"),
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/shorten", h.handleShorten)
	mux.HandleFunc("/", h.handleRedirect)
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
	LongURL   string `json:"long_url"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	defer r.Body.Close()

	var req shortenRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	longURL, err := normalizeLongURL(req.URL)
	if err != nil {
		h.writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	shortCode, err := h.shortener.Shorten(r.Context(), longURL)
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "failed to shorten url")
		return
	}

	h.writeJSON(w, http.StatusCreated, shortenResponse{
		ShortCode: string(shortCode),
		ShortURL:  h.buildShortURL(r, shortCode),
		LongURL:   string(longURL),
	})
}

func (h *Handler) handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if r.URL.Path == "/" {
		h.writeJSONError(w, http.StatusNotFound, "short url not found")
		return
	}

	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	if shortCode == "" || strings.Contains(shortCode, "/") {
		h.writeJSONError(w, http.StatusNotFound, "short url not found")
		return
	}

	longURL, err := h.resolver.Resolve(r.Context(), domain.ShortURL(shortCode))
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "failed to resolve short url")
		return
	}
	if longURL == "" {
		h.writeJSONError(w, http.StatusNotFound, "short url not found")
		return
	}

	http.Redirect(w, r, string(longURL), http.StatusFound)
}

func normalizeLongURL(raw string) (domain.LongURL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("url is required")
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("invalid url")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("url must use http or https")
	}
	if parsed.Host == "" {
		return "", fmt.Errorf("url host is required")
	}

	return domain.LongURL(parsed.String()), nil
}

func (h *Handler) buildShortURL(r *http.Request, shortCode domain.ShortURL) string {
	if h.baseURL != "" {
		return h.baseURL + "/" + string(shortCode)
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if forwardedProto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); forwardedProto != "" {
		scheme = forwardedProto
	}

	return fmt.Sprintf("%s://%s/%s", scheme, r.Host, shortCode)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) writeJSONError(w http.ResponseWriter, status int, msg string) {
	h.writeJSON(w, status, errorResponse{Error: msg})
}
