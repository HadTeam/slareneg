package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server/internal/config"
	"server/internal/wire"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	setupLogging(cfg.Logging)
	slog.Info("starting slareneg game server")

	app, err := wire.InitializeApplication(cfg)
	if err != nil {
		slog.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if err := startServices(app, cfg); err != nil {
		slog.Error("failed to start services", "error", err)
		os.Exit(1)
	}

	select {
	case sig := <-sigChan:
		slog.Info("received signal, shutting down", "signal", sig)
	case <-ctx.Done():
		slog.Info("context cancelled, shutting down")
	}

	slog.Info("shutting down server")
	stopServices(app)
	slog.Info("server shutdown complete")
}

func startServices(app *wire.Application, cfg *config.Config) error {
	go func() {
		if err := app.Lobby.Start(); err != nil {
			slog.Error("lobby service failed to start", "error", err)
		}
	}()

	setupHTTPRoutes(app)

	go func() {
		slog.Info("starting HTTP server", "addr", cfg.GetServerAddr())

		server := &http.Server{
			Addr:         cfg.GetServerAddr(),
			ReadTimeout:  time.Duration(cfg.Server.ReadTimeout),
			WriteTimeout: time.Duration(cfg.Server.WriteTimeout),
		}

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server failed", "error", err)
		}
	}()

	return nil
}

func setupHTTPRoutes(app *wire.Application) {
	http.HandleFunc("/api/auth/register", app.AuthService.RegisterHandler)
	http.HandleFunc("/api/auth/login", app.AuthService.LoginHandler)
	http.HandleFunc("/api/game/ws", app.AuthService.AuthMiddleware(app.WSServer.HandleWebSocket))
	http.HandleFunc("/health", healthCheckHandler(app))
	http.HandleFunc("/api/cache/stats", app.AuthService.AuthMiddleware(cacheStatsHandler(app)))

	staticDir := app.Config.Server.StaticDir
	if _, err := os.Stat(staticDir); err == nil {
		http.Handle("/", http.FileServer(http.Dir(staticDir)))
		slog.Info("serving static files", "dir", staticDir)
	}
}

func healthCheckHandler(app *wire.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		stats := app.Cache.GetCacheStats()
		health := map[string]interface{}{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"cache":     stats,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(health); err != nil {
			slog.Error("failed to encode health response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func cacheStatsHandler(app *wire.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		stats := app.Cache.GetCacheStats()
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			slog.Error("failed to encode cache stats", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func stopServices(app *wire.Application) {
	if err := app.Lobby.Stop(); err != nil {
		slog.Error("error stopping lobby service", "error", err)
	}

	if err := app.WSServer.StopServer(); err != nil {
		slog.Error("error stopping websocket server", "error", err)
	}
}

func setupLogging(cfg config.LoggingConfig) {
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
