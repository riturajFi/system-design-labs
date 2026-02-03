package httptracker

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"notification-system/internal/core/model"
)

type Tracker struct {
	endpoint string
	client   *http.Client
}

func New(endpoint string) *Tracker {
	return &Tracker{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 500 * time.Millisecond,
		},
	}
}

func (t *Tracker) Track(
	_ context.Context,
	e model.NotificationEvent,
) error {

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		t.endpoint,
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	_, err = t.client.Do(req)
	return err
}
