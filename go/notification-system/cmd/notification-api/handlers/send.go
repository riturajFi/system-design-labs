package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"notification-system/internal/core/engine"
	"notification-system/internal/core/model"
)

type SendRequest struct {
	EventID    string            `json:"event_id"`
	UserID     int64             `json:"user_id"`
	Channel    model.Channel     `json:"channel"`
	TemplateID string            `json:"template_id"`
	Params     map[string]string `json:"params"`
}

func SendHandler(e *engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: What id this?
		defer r.Body.Close()

		var req SendRequest

		// TODO : What is this?
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		n := model.Notification{
			EventID:    req.EventID,
			UserID:     req.UserID,
			Channel:    req.Channel,
			TemplateID: req.TemplateID,
			Params:     req.Params,
			CreatedAt:  time.Now().Unix(),
		}

		appKey := r.Header.Get("X-App-Key")
		appSecret := r.Header.Get("X-App-Secret")

		err := e.Send(context.Background(), appKey, appSecret, n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("enqueued"))
	}
}
