package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	addr := ":8080"

	mux := http.NewServeMux()

	// Health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start server
	go func() {
		fmt.Println("api: listening on", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("api: listen error:", err)
			os.Exit(1)
		}
	}()

	// Wait for SIGINT / SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	fmt.Println("api: shutdown start")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("api: shutdown error:", err)
	}

	fmt.Println("api: shutdown complete")
}
