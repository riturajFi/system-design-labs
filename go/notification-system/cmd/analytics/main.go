package main

import (
	"net/http"

	"notification-system/internal/analytics"
)

func main() {
	h := analytics.NewHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("/ingest", h.Ingest)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	http.ListenAndServe(":8090", mux)
}
