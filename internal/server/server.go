package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/TelitsynNikita/test-example-for-maksec/internal/api/handler"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/api/middleware"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/config"
	"github.com/TelitsynNikita/test-example-for-maksec/internal/service"

	_ "github.com/TelitsynNikita/test-example-for-maksec/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	CreateServer   *http.Server
	CallbackServer *http.Server
}

func New(services *service.Services, cfg *config.Config, logger *slog.Logger) *Server {
	createHandler := handler.NewCreateHandler(services.ScriptService)
	callbackHandler := handler.NewCallbackHandler(services.EventService)
	healthHandler := handler.NewHealthHandler()

	createMux := http.NewServeMux()
	createMux.HandleFunc("POST /create", createHandler.Handle)
	createMux.HandleFunc("GET /health", healthHandler.Handle)
	createMux.HandleFunc("GET /health/ready", healthHandler.Handle)

	createMux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	callbackMux := http.NewServeMux()
	callbackMux.Handle("POST /callback", middleware.Auth(cfg.Callback.Token)(
		http.HandlerFunc(callbackHandler.Handle),
	))

	middlewareChain := middleware.Chain(
		middleware.MaxBodySize(1024*1024),
		middleware.RateLimit(10, 20),
		middleware.Logging(logger),
	)

	return &Server{
		CreateServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.CreatePort),
			Handler:      middlewareChain(createMux),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
		CallbackServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Server.CallbackPort),
			Handler:      middlewareChain(callbackMux),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
	}
}

func (s *Server) Start(logger *slog.Logger) {
	go func() {
		logger.Info("create API server starting", "port", s.CreateServer.Addr)
		if err := s.CreateServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("create server failed", "error", err)
		}
	}()

	go func() {
		logger.Info("callback API server starting", "port", s.CallbackServer.Addr)
		if err := s.CallbackServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("callback server failed", "error", err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context, logger *slog.Logger) {
	logger.Info("shutting down servers...")

	if err := s.CreateServer.Shutdown(ctx); err != nil {
		logger.Error("create server forced to shutdown", "error", err)
	}

	if err := s.CallbackServer.Shutdown(ctx); err != nil {
		logger.Error("callback server forced to shutdown", "error", err)
	}

	logger.Info("servers stopped")
}
