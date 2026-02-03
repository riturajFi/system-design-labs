package main

import (
	"net/http"

	"notification-system/internal/config"
	"notification-system/internal/core/engine"
	"notification-system/internal/core/model"
	authstatic "notification-system/internal/modules/auth/static"
	contactmemory "notification-system/internal/modules/contact/memory"
	ratelimitfixed "notification-system/internal/modules/ratelimit/fixed"
	settingsstatic "notification-system/internal/modules/settings/static"
	"notification-system/internal/observability/logging"
	"notification-system/internal/observability/metrics"
)

func main() {
	cfg := config.Load()
	logger := logging.New(cfg.ServiceName, cfg.Env)

	registry := metrics.NewRegistry()
	metricsHandler := metrics.NewHandler(registry)

	auth := authstatic.New(map[string]string{
		"demo-app": "demo-secret",
	})

	settings := settingsstatic.New(map[settingsstatic.Key]bool{
		{UserID: 1, Channel: "email"}: true,
		{UserID: 1, Channel: "sms"}:   false,
	})

	contactStore := contactmemory.New(
		[]model.User{
			{ID: 1, Email: "user1@example.com", PhoneNumber: "+1234567890"},
		},
		[]model.Device{
			{ID: 1, UserID: 1, Token: "ios-token-1", Platform: model.ChannelPushIOS},
			{ID: 2, UserID: 1, Token: "android-token-1", Platform: model.ChannelPushAndroid},
		},
	)
	_ = contactStore

	rateLimiter := ratelimitfixed.New(5, registry) // 5 per user/channel/min

	e := engine.New(cfg, logger, registry, engine.Deps{
		Auth:      auth,
		RateLimit: rateLimiter,
		Settings:  settings,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/health", e.Health)
	mux.Handle("/metrics", metricsHandler)

	logger.Info("starting http server on :" + cfg.HTTPPort)
	http.ListenAndServe(":"+cfg.HTTPPort, mux)
}
