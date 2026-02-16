package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Server struct {
	router *chi.Mux
	logger *zap.Logger
	port   int
}

func NewServer(port int, logger *zap.Logger) *Server {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	return &Server{
		router: r,
		logger: logger,
		port:   port,
	}
}

func (s *Server) Router() *chi.Mux {
	return s.router
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	s.logger.Info("Starting HTTP server", zap.String("address", addr))

	if err := http.ListenAndServe(addr, s.router); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")
	return nil
}
