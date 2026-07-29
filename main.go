package main

import (
	"context"
	"edge-gateway/internal/app"
	"edge-gateway/internal/config"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// --- Config & Logger ---
	cfg := config.Load()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// --- Server Setup ---
	srv := app.New(cfg, log)
	srv.StartWatchers()

	// --- Start HTTP Server ---
	httpSrv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      srv.Handler(),
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	}

	go func() {
		log.Info("Gateway starting", slog.String("port", cfg.App.Port), slog.String("dir", cfg.App.Dir))
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// --- Graceful Shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down HTTP server...")
	srv.StopWatchers()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Error("Server forced to shutdown", slog.String("error", err.Error()))
	}

	log.Info("HTTP server shutdown successfully")
	os.Exit(0)
}
