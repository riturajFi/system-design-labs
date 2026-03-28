package storagehttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"urlshortener/internal/domain"
)

type repository interface {
	GetByShortURL(ctx context.Context, shortURL domain.ShortURL) (*domain.LongURL, error)
	GetByLongURL(ctx context.Context, longURL domain.LongURL) (*domain.ShortURL, error)
	Save(ctx context.Context, id domain.ID, shortURL domain.ShortURL, longURL domain.LongURL) error
}

type idAllocator interface {
	NextID(ctx context.Context) (domain.ID, error)
}

type Handler struct {
	repo  repository
	idgen idAllocator
}

type mappingPayload struct {
	ID       domain.ID       `json:"id"`
	ShortURL domain.ShortURL `json:"short_url"`
	LongURL  domain.LongURL  `json:"long_url"`
}

type idResponse struct {
	ID domain.ID `json:"id"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler(repo repository, idgen idAllocator) *Handler {
	return &Handler{
		repo:  repo,
		idgen: idgen,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/ids/next", h.handleNextID)
	mux.HandleFunc("/mappings/short/", h.handleGetByShort)
	mux.HandleFunc("/mappings/long", h.handleGetByLong)
	mux.HandleFunc("/mappings", h.handleSave)
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleNextID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id, err := h.idgen.NextID(r.Context())
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "failed to allocate id")
		return
	}

	h.writeJSON(w, http.StatusOK, idResponse{ID: id})
}

func (h *Handler) handleGetByShort(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	shortPart := strings.TrimPrefix(r.URL.Path, "/mappings/short/")
	shortValue, err := url.PathUnescape(shortPart)
	if err != nil || shortValue == "" {
		h.writeJSONError(w, http.StatusBadRequest, "invalid short url")
		return
	}

	longURL, err := h.repo.GetByShortURL(r.Context(), domain.ShortURL(shortValue))
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "failed to read mapping")
		return
	}
	if longURL == nil {
		h.writeJSONError(w, http.StatusNotFound, "mapping not found")
		return
	}

	shortURL, err := h.repo.GetByLongURL(r.Context(), *longURL)
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "failed to read mapping")
		return
	}
	if shortURL == nil {
		h.writeJSONError(w, http.StatusInternalServerError, "mapping lookup inconsistent")
		return
	}

	h.writeJSON(w, http.StatusOK, mappingPayload{
		ShortURL: *shortURL,
		LongURL:  *longURL,
	})
}

func (h *Handler) handleGetByLong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	longValue := strings.TrimSpace(r.URL.Query().Get("url"))
	if longValue == "" {
		h.writeJSONError(w, http.StatusBadRequest, "url query parameter is required")
		return
	}

	shortURL, err := h.repo.GetByLongURL(r.Context(), domain.LongURL(longValue))
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "failed to read mapping")
		return
	}
	if shortURL == nil {
		h.writeJSONError(w, http.StatusNotFound, "mapping not found")
		return
	}

	longURL, err := h.repo.GetByShortURL(r.Context(), *shortURL)
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "failed to read mapping")
		return
	}
	if longURL == nil {
		h.writeJSONError(w, http.StatusInternalServerError, "mapping lookup inconsistent")
		return
	}

	h.writeJSON(w, http.StatusOK, mappingPayload{
		ShortURL: *shortURL,
		LongURL:  *longURL,
	})
}

func (h *Handler) handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	defer r.Body.Close()

	var payload mappingPayload
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if payload.ShortURL == "" || payload.LongURL == "" {
		h.writeJSONError(w, http.StatusBadRequest, "short_url and long_url are required")
		return
	}

	if err := h.repo.Save(r.Context(), payload.ID, payload.ShortURL, payload.LongURL); err != nil {
		h.writeJSONError(w, http.StatusConflict, fmt.Sprintf("save mapping: %v", err))
		return
	}

	h.writeJSON(w, http.StatusCreated, payload)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *Handler) writeJSONError(w http.ResponseWriter, status int, msg string) {
	h.writeJSON(w, status, errorResponse{Error: msg})
}
