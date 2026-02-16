package handlers

import (
	"net/http"

	"github.com/akhilbidhuri/SimplNews/internal/repository/postgres"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	articleRepo *postgres.ArticleRepository
	logger      *zap.Logger
}

func NewCategoryHandler(articleRepo *postgres.ArticleRepository, logger *zap.Logger) *CategoryHandler {
	return &CategoryHandler{
		articleRepo: articleRepo,
		logger:      logger,
	}
}

// HandleCategory returns articles filtered by category
func (h *CategoryHandler) HandleCategory(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		writeError(w, http.StatusBadRequest, "category parameter is required")
		return
	}

	limit := parseLimit(r.URL.Query().Get("limit"), 5, 20)

	articles, err := h.articleRepo.FindByCategory(category, limit)
	if err != nil {
		h.logger.Error("Failed to find articles by category", zap.Error(err), zap.String("category", category))
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
