package storageclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"urlshortener/internal/domain"
	filerepo "urlshortener/internal/storage/repository/file"
	"urlshortener/internal/storagehttp"
)

func TestClientRoundTripPersistsToFileBackedStorage(t *testing.T) {
	storageFile := filepath.Join(t.TempDir(), "mappings.json")

	repo, err := filerepo.Open(storageFile)
	if err != nil {
		t.Fatalf("open file repository: %v", err)
	}

	handler := storagehttp.NewHandler(repo, repo)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := httptest.NewServer(mux)
	defer server.Close()

	client := New(server.URL)
	ctx := context.Background()

	id, err := client.NextID(ctx)
	if err != nil {
		t.Fatalf("next id: %v", err)
	}
	if id != 1 {
		t.Fatalf("unexpected id: got %d want %d", id, 1)
	}

	err = client.Save(ctx, Mapping{
		ID:       id,
		ShortURL: domain.ShortURL("1"),
		LongURL:  domain.LongURL("https://example.com/stored"),
	})
	if err != nil {
		t.Fatalf("save mapping: %v", err)
	}

	byShort, err := client.GetByShortURL(ctx, domain.ShortURL("1"))
	if err != nil {
		t.Fatalf("get by short url: %v", err)
	}
	if byShort == nil || byShort.LongURL != domain.LongURL("https://example.com/stored") {
		t.Fatalf("unexpected mapping by short: %#v", byShort)
	}

	reloadedRepo, err := filerepo.Open(storageFile)
	if err != nil {
		t.Fatalf("reopen file repository: %v", err)
	}

	longURL, err := reloadedRepo.GetByShortURL(ctx, domain.ShortURL("1"))
	if err != nil {
		t.Fatalf("reloaded get by short url: %v", err)
	}
	if longURL == nil || *longURL != domain.LongURL("https://example.com/stored") {
		t.Fatalf("unexpected persisted mapping: %#v", longURL)
	}

	nextID, err := reloadedRepo.NextID(ctx)
	if err != nil {
		t.Fatalf("reloaded next id: %v", err)
	}
	if nextID != 2 {
		t.Fatalf("unexpected next id after reload: got %d want %d", nextID, 2)
	}
}
