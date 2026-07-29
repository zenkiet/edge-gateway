package app

import (
	"edge-gateway/internal/config"
	"edge-gateway/internal/pkg/middleware"
	"edge-gateway/internal/proxy"
	"edge-gateway/internal/static"
	"log/slog"
	"net/http"
)

type Server struct {
	manager *static.Manager
	logger  *slog.Logger
	mux     *http.ServeMux
}

func New(cfg *config.Config, logger *slog.Logger) *Server {
	mux := http.NewServeMux()

	// Initialize Static Domain
	manager := static.NewManager(cfg.App.Dir, logger)
	staticHandler := static.NewHandler(manager)

	// Reverse Proxy
	if cfg.App.UpstreamURL != "" {
		proxyHandler := proxy.NewProxy(cfg.App.UpstreamURL, manager.ResolveBundleID, logger)
		if proxyHandler != nil {
			logger.Info("Proxy enabled", "upstream", cfg.App.UpstreamURL)
			mux.Handle("/api/", proxyHandler)
		}
	}

	mux.Handle("/", staticHandler)

	return &Server{
		mux:     mux,
		manager: manager,
		logger:  logger,
	}

}

func (s *Server) StartWatchers() {
	s.manager.StartWatcher()
}

func (s *Server) StopWatchers() {
	s.manager.StopWatcher()
}

func (s *Server) Handler() http.Handler {
	handler := middleware.RequestLogger(s.logger)(s.mux)
	return middleware.Recovery(s.logger)(handler)
}
