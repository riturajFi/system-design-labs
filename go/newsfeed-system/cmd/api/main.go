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

/*
loggingMiddleware:
- logs method, path, status, latency
- no external deps
*/
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// capture status code
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(lrw, r)

		latencyMs := time.Since(start).Milliseconds()
		fmt.Printf(
			"req method=%s path=%s status=%d latency_ms=%d\n",
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			latencyMs,
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (l *loggingResponseWriter) WriteHeader(code int) {
	l.statusCode = code
	l.ResponseWriter.WriteHeader(code)
}

func main() {
	addr := ":8080"

	mux := http.NewServeMux()

	// health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// wrap mux with logging
	handler := loggingMiddleware(mux)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// start server
	go func() {
		fmt.Println("api: listening on", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("api: listen error:", err)
			os.Exit(1)
		}
	}()

	// shutdown handling
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
