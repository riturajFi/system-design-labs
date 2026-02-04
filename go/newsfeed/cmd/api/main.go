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
	"newsfeed/internal/ratelimit"
)

/* ---------- logging ---------- */

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: 200}
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

/* ---------- auth ---------- */

type ctxUserIDKey struct{}

func authMiddleware(a auth.Authenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		h := r.Header.Get("Authorization")
		if h == "" {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}

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

		ctx := context.WithValue(r.Context(), ctxUserIDKey{}, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

/* ---------- rate limit ---------- */

func rateLimitMiddleware(limiter ratelimit.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		userID := r.Context().Value(ctxUserIDKey{}).(string)

		if !limiter.Allow(r.Context(), userID) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	addr := ":8080"

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/whoami", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Context().Value(ctxUserIDKey{}).(string)))
	})

	authenticator := &auth.MockAuthenticator{}
	limiter := ratelimit.NewMemoryLimiter(3, 10*time.Second)

	handler :=
		loggingMiddleware(
			authMiddleware(
				authenticator,
				rateLimitMiddleware(limiter, mux),
			),
		)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		fmt.Println("api: listening on", addr)
		server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	server.Shutdown(context.Background())
}
