package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"urlshortener/internal/idgen/memory"
	"urlshortener/internal/resolution"
	"urlshortener/internal/shortening"
	cachememory "urlshortener/internal/storage/cache/memory"
	repomemory "urlshortener/internal/storage/repository/memory"
)

func TestShortenAndRedirect(t *testing.T) {
	repo := repomemory.New()
	cache := cachememory.New()
	idg := memory.New(0)
	shortener := shortening.NewService(repo, idg)
	resolver := resolution.NewService(cache, repo)
	handler := NewHandler(shortener, resolver, "")

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	body := bytes.NewBufferString(`{"url":"https://example.com/some/path?q=1"}`)
	resp, err := http.Post(server.URL+"/shorten", "application/json", body)
	if err != nil {
		t.Fatalf("shorten request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected shorten status: got %d want %d", resp.StatusCode, http.StatusCreated)
	}

	var created shortenResponse
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode shorten response: %v", err)
	}

	if created.ShortCode != "1" {
		t.Fatalf("unexpected short code: got %q want %q", created.ShortCode, "1")
	}

	client := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	redirectResp, err := client.Get(server.URL + "/" + created.ShortCode)
	if err != nil {
		t.Fatalf("redirect request failed: %v", err)
	}
	defer redirectResp.Body.Close()

	if redirectResp.StatusCode != http.StatusFound {
		t.Fatalf("unexpected redirect status: got %d want %d", redirectResp.StatusCode, http.StatusFound)
	}

	if location := redirectResp.Header.Get("Location"); location != "https://example.com/some/path?q=1" {
		t.Fatalf("unexpected redirect location: got %q", location)
	}
}

func TestShortenIsIdempotent(t *testing.T) {
	repo := repomemory.New()
	cache := cachememory.New()
	idg := memory.New(0)
	shortener := shortening.NewService(repo, idg)
	resolver := resolution.NewService(cache, repo)
	handler := NewHandler(shortener, resolver, "")

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	first := shortenURL(t, server.URL, "https://example.com/repeat")
	second := shortenURL(t, server.URL, "https://example.com/repeat")

	if first.ShortCode != second.ShortCode {
		t.Fatalf("expected idempotent shorten result, got %q and %q", first.ShortCode, second.ShortCode)
	}
}

func shortenURL(t *testing.T, baseURL string, longURL string) shortenResponse {
	t.Helper()

	body := bytes.NewBufferString(`{"url":"` + longURL + `"}`)
	resp, err := http.Post(baseURL+"/shorten", "application/json", body)
	if err != nil {
		t.Fatalf("shorten request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected shorten status: got %d want %d", resp.StatusCode, http.StatusCreated)
	}

	var result shortenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode shorten response: %v", err)
	}

	return result
}
