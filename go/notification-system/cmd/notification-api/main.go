package main

import (
	"log"
	"net/http"

	"notification-system/internal/core/engine"
)

func main() {
	e := engine.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", e.Health)

	log.Println("notification-api listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
