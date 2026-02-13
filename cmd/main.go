package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/akhilbidhuri/SimplNews/internal/pkg/config"
	"github.com/akhilbidhuri/SimplNews/internal/pkg/logger"
	"github.com/akhilbidhuri/SimplNews/internal/repository/postgres"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(cfg.Logging.Level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Infow("SimplNews API starting",
		"version", "1.0.0",
		"port", cfg.Server.Port,
		"environment", "development",
	)

	// Connect to database
	db, err := postgres.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalw("Failed to connect to database", "error", err)
	}
	defer db.Close()

	log.Infow("Connected to PostgreSQL", "database", cfg.Database.Name)

	// TODO: Initialize repositories, services, handlers, and HTTP server
	// For now, just start a simple health check server
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","message":"SimplNews API is running"}`)
	})

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        mux,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Infow("Starting HTTP server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalw("Server error", "error", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Infow("Received shutdown signal", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorw("Server shutdown error", "error", err)
	}

	log.Infow("SimplNews API stopped")
}
