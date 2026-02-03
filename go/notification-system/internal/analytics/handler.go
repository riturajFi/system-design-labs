package analytics

import (
	"encoding/json"
	"net/http"
	"sync"

	"notification-system/internal/core/model"
)

type Handler struct {
	mu     sync.Mutex
	events []model.NotificationEvent
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Ingest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var e model.NotificationEvent
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "invalid event", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	h.events = append(h.events, e)
	h.mu.Unlock()

	w.WriteHeader(http.StatusAccepted)
}
