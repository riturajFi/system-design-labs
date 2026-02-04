package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"newsfeed/internal/auth"
)

/*
logging middleware (unchanged)
*/
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(lrw, r)

		fmt.Printf(
			"req method=%s path=%s status=%d latency_ms=%d\n",
			r.Method,
			r.URL.Path,
			lrw.statusCode,
			time.Since(start).Milliseconds(),
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

/*
auth middleware
*/
func authMiddleware(a auth.Authenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow health without auth
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		h := r.Header.Get("Authorization")
		if h == "" {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}

		// expect: Authorization: Bearer token-<userID>
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid auth header", http.StatusUnauthorized)
			return
		}

		userID, ok := a.Authenticate(r.Context(), parts[1])
		if !ok {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// attach userID to context
		ctx := context.WithValue(r.Context(), ctxUserIDKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type ctxUserIDKey struct{}

func main() {
	addr := ":8080"

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// example protected endpoint (temporary)
	mux.HandleFunc("/whoami", func(w http.ResponseWriter, r *http.Request) {
		uid := r.Context().Value(ctxUserIDKey{}).(string)
		w.Write([]byte(uid))
	})

	authenticator := &auth.MockAuthenticator{}

	handler := loggingMiddleware(
		authMiddleware(authenticator, mux),
	)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		fmt.Println("api: listening on", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("api: listen error:", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
