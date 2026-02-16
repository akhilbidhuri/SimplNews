package handlers

import (
	"net/http"

	"github.com/akhilbidhuri/SimplNews/internal/repository/postgres"
	"go.uber.org/zap"
)

type SearchHandler struct {
	articleRepo *postgres.ArticleRepository
	logger      *zap.Logger
}

func NewSearchHandler(articleRepo *postgres.ArticleRepository, logger *zap.Logger) *SearchHandler {
	return &SearchHandler{
		articleRepo: articleRepo,
		logger:      logger,
	}
}

// HandleSearch performs full-text search on articles
func (h *SearchHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		writeError(w, http.StatusBadRequest, "query parameter is required")
		return
	}

	limit := parseLimit(r.URL.Query().Get("limit"), 5, 20)

	articles, err := h.articleRepo.SearchByText(query, limit)
	if err != nil {
		h.logger.Error("Failed to search articles", zap.Error(err), zap.String("query", query))
		writeError(w, http.StatusInternalServerError, "Failed to search articles")
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"count":  len(articles),
		"data":   articles,
	}

	writeJSON(w, response)
}
