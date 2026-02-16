package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/akhilbidhuri/SimplNews/internal/domain/services"
	"go.uber.org/zap"
)

type QueryHandler struct {
	queryService *services.QueryService
	logger       *zap.Logger
}

func NewQueryHandler(queryService *services.QueryService, logger *zap.Logger) *QueryHandler {
	return &QueryHandler{
		queryService: queryService,
		logger:       logger,
	}
}

type QueryRequest struct {
	Query        string    `json:"query"`
	UserLocation *Location `json:"user_location,omitempty"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// HandleQuery processes natural language queries with LLM-powered intent extraction
func (h *QueryHandler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request", zap.Error(err))
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Query == "" {
		writeError(w, http.StatusBadRequest, "Query parameter is required")
		return
	}

	limit := parseLimit(r.URL.Query().Get("limit"), 5, 20)

	var userLat, userLon *float64
	if req.UserLocation != nil {
		userLat = &req.UserLocation.Latitude
		userLon = &req.UserLocation.Longitude
	}

	articles, intentData, err := h.queryService.ProcessQuery(r.Context(), req.Query, userLat, userLon, limit)
	if err != nil {
		h.logger.Error("Failed to process query", zap.Error(err), zap.String("query", req.Query))
		writeError(w, http.StatusInternalServerError, "Failed to process query")
		return
	}

	h.logger.Info("Query processed",
		zap.String("query", req.Query),
		zap.String("intent", intentData.Intent),
		zap.Int("results", len(articles)),
	)

	response := map[string]interface{}{
		"status": "success",
		"count":  len(articles),
		"intent": intentData.Intent,
		"data":   articles,
	}

	writeJSON(w, response)
}
