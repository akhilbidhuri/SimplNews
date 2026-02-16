package handlers

import (
	"net/http"

	"github.com/akhilbidhuri/SimplNews/internal/repository/postgres"
	"go.uber.org/zap"
)

type ScoreHandler struct {
	articleRepo *postgres.ArticleRepository
	logger      *zap.Logger
}

func NewScoreHandler(articleRepo *postgres.ArticleRepository, logger *zap.Logger) *ScoreHandler {
	return &ScoreHandler{
		articleRepo: articleRepo,
		logger:      logger,
	}
}

// HandleScore returns articles with relevance score >= min_score
func (h *ScoreHandler) HandleScore(w http.ResponseWriter, r *http.Request) {
	minScoreStr := r.URL.Query().Get("min_score")
	if minScoreStr == "" {
		writeError(w, http.StatusBadRequest, "min_score parameter is required")
		return
	}

	minScore := parseFloat(minScoreStr, 0.0)
	if minScore < 0 || minScore > 1 {
		writeError(w, http.StatusBadRequest, "min_score must be between 0 and 1")
		return
	}

	limit := parseLimit(r.URL.Query().Get("limit"), 5, 20)

	articles, err := h.articleRepo.FindByMinScore(minScore, limit)
	if err != nil {
		h.logger.Error("Failed to find articles by score", zap.Error(err), zap.Float64("min_score", minScore))
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
