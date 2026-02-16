package api

import (
	"encoding/json"
	"net/http"

	"github.com/akhilbidhuri/SimplNews/internal/api/handlers"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(
	r *chi.Mux,
	queryHandler *handlers.QueryHandler,
	categoryHandler *handlers.CategoryHandler,
	scoreHandler *handlers.ScoreHandler,
	searchHandler *handlers.SearchHandler,
	sourceHandler *handlers.SourceHandler,
	nearbyHandler *handlers.NearbyHandler,
) {
	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "healthy",
		})
	})

	// API v1 routes
	r.Route("/api/v1/news", func(r chi.Router) {
		// Main LLM-powered query endpoint
		r.Post("/query", queryHandler.HandleQuery)

		// Direct filtering endpoints
		r.Get("/category", categoryHandler.HandleCategory)
		r.Get("/score", scoreHandler.HandleScore)
		r.Get("/search", searchHandler.HandleSearch)
		r.Get("/source", sourceHandler.HandleSource)
		r.Get("/nearby", nearbyHandler.HandleNearby)
	})
}
