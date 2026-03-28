package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"urlshortener/internal/domain"
	"urlshortener/internal/storage/repository"
)

type record struct {
	ID       domain.ID       `json:"id"`
	ShortURL domain.ShortURL `json:"short_url"`
	LongURL  domain.LongURL  `json:"long_url"`
}

type persistedState struct {
	NextID   domain.ID `json:"next_id"`
	Mappings []record  `json:"mappings"`
}

type Repository struct {
	mu        sync.RWMutex
	filePath  string
	nextID    domain.ID
	byShort   map[domain.ShortURL]record
	byLongURL map[domain.LongURL]record
}

func Open(filePath string) (*Repository, error) {
	repo := &Repository{
		filePath:  filePath,
		byShort:   make(map[domain.ShortURL]record),
		byLongURL: make(map[domain.LongURL]record),
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) GetByShortURL(_ context.Context, shortURL domain.ShortURL) (*domain.LongURL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rec, ok := r.byShort[shortURL]
	if !ok {
		return nil, nil
	}

	result := rec.LongURL
	return &result, nil
}

func (r *Repository) GetByLongURL(_ context.Context, longURL domain.LongURL) (*domain.ShortURL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rec, ok := r.byLongURL[longURL]
	if !ok {
		return nil, nil
	}

	result := rec.ShortURL
	return &result, nil
}

func (r *Repository) Save(_ context.Context, id domain.ID, shortURL domain.ShortURL, longURL domain.LongURL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byShort[shortURL]; exists {
		return fmt.Errorf("short url already exists: %s", shortURL)
	}
	if _, exists := r.byLongURL[longURL]; exists {
		return fmt.Errorf("long url already exists: %s", longURL)
	}

	rec := record{
		ID:       id,
		ShortURL: shortURL,
		LongURL:  longURL,
	}

	r.byShort[shortURL] = rec
	r.byLongURL[longURL] = rec

	if err := r.persistLocked(); err != nil {
		delete(r.byShort, shortURL)
		delete(r.byLongURL, longURL)
		return err
	}

	return nil
}

func (r *Repository) NextID(_ context.Context) (domain.ID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nextID++

	if err := r.persistLocked(); err != nil {
		r.nextID--
		return 0, err
	}

	return r.nextID, nil
}

func (r *Repository) load() error {
	if r.filePath == "" {
		return fmt.Errorf("storage file path is required")
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read storage file: %w", err)
	}

	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("unmarshal storage file: %w", err)
	}

	r.nextID = state.NextID
	for _, rec := range state.Mappings {
		r.byShort[rec.ShortURL] = rec
		r.byLongURL[rec.LongURL] = rec
	}

	return nil
}

func (r *Repository) persistLocked() error {
	if err := os.MkdirAll(filepath.Dir(r.filePath), 0o755); err != nil {
		return fmt.Errorf("create storage directory: %w", err)
	}

	state := persistedState{
		NextID:   r.nextID,
		Mappings: make([]record, 0, len(r.byShort)),
	}
	for _, rec := range r.byShort {
		state.Mappings = append(state.Mappings, rec)
	}

	sort.Slice(state.Mappings, func(i, j int) bool {
		return state.Mappings[i].ID < state.Mappings[j].ID
	})

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal storage state: %w", err)
	}

	tmpPath := r.filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return fmt.Errorf("write temp storage file: %w", err)
	}
	if err := os.Rename(tmpPath, r.filePath); err != nil {
		return fmt.Errorf("rename storage file: %w", err)
	}

	return nil
}

var _ repository.URLRepository = (*Repository)(nil)
