package handlers

import (
	"net/http"

	"github.com/akhilbidhuri/SimplNews/internal/repository/postgres"
	"go.uber.org/zap"
)

type NearbyHandler struct {
	articleRepo *postgres.ArticleRepository
	logger      *zap.Logger
}

func NewNearbyHandler(articleRepo *postgres.ArticleRepository, logger *zap.Logger) *NearbyHandler {
	return &NearbyHandler{
		articleRepo: articleRepo,
		logger:      logger,
	}
}

// HandleNearby returns articles near a specific location
func (h *NearbyHandler) HandleNearby(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("latitude")
	lonStr := r.URL.Query().Get("longitude")

	if latStr == "" || lonStr == "" {
		writeError(w, http.StatusBadRequest, "latitude and longitude parameters are required")
		return
	}

	lat := parseFloat(latStr, 0)
	lon := parseFloat(lonStr, 0)

	// Validate coordinates
	if lat < -90 || lat > 90 {
		writeError(w, http.StatusBadRequest, "latitude must be between -90 and 90")
		return
	}
	if lon < -180 || lon > 180 {
		writeError(w, http.StatusBadRequest, "longitude must be between -180 and 180")
		return
	}

	radiusKm := parseFloat(r.URL.Query().Get("radius_km"), 50.0)
	if radiusKm > 500 {
		radiusKm = 500 // Max 500km
	}

	limit := parseLimit(r.URL.Query().Get("limit"), 5, 20)

	articles, err := h.articleRepo.FindNearby(lat, lon, radiusKm, limit)
	if err != nil {
		h.logger.Error("Failed to find nearby articles",
			zap.Error(err),
			zap.Float64("latitude", lat),
			zap.Float64("longitude", lon),
			zap.Float64("radius_km", radiusKm),
		)
		writeError(w, http.StatusInternalServerError, "Failed to retrieve articles")
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"count":  len(articles),
		"data":   articles,
	}

	writeJSON(w, response)
}
