package handlers

import (
	"net/http"

	"github.com/akhilbidhuri/SimplNews/internal/repository/postgres"
	"go.uber.org/zap"
)

type SourceHandler struct {
	articleRepo *postgres.ArticleRepository
	logger      *zap.Logger
}

func NewSourceHandler(articleRepo *postgres.ArticleRepository, logger *zap.Logger) *SourceHandler {
	return &SourceHandler{
		articleRepo: articleRepo,
		logger:      logger,
	}
}

// HandleSource returns articles from a specific news source
func (h *SourceHandler) HandleSource(w http.ResponseWriter, r *http.Request) {
	source := r.URL.Query().Get("source")
	if source == "" {
		writeError(w, http.StatusBadRequest, "source parameter is required")
		return
	}

	limit := parseLimit(r.URL.Query().Get("limit"), 5, 20)

	articles, err := h.articleRepo.FindBySource(source, limit)
	if err != nil {
		h.logger.Error("Failed to find articles by source", zap.Error(err), zap.String("source", source))
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
