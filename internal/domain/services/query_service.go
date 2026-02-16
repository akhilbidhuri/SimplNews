package services

import (
	"context"
	"fmt"

	"github.com/akhilbidhuri/SimplNews/internal/domain/models"
	"github.com/akhilbidhuri/SimplNews/internal/repository/postgres"
)

type QueryService struct {
	articleRepo *postgres.ArticleRepository
	llmService  *LLMService
}

func NewQueryService(articleRepo *postgres.ArticleRepository, llmService *LLMService) *QueryService {
	return &QueryService{
		articleRepo: articleRepo,
		llmService:  llmService,
	}
}

// ProcessQuery analyzes user query with LLM and routes to appropriate endpoint
func (s *QueryService) ProcessQuery(ctx context.Context, query string, userLat, userLon *float64, limit int) ([]models.Article, *IntentData, error) {
	// Extract intent and entities using LLM
	intentData, err := s.llmService.ExtractIntent(ctx, query, userLat, userLon)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to extract intent: %w", err)
	}

	var articles []models.Article

	// Route based on detected intent
	switch intentData.Intent {
	case "nearby":
		// Geospatial search - requires user location
		if userLat != nil && userLon != nil {
			radiusKm := 50.0 // Default 50km radius
			articles, err = s.articleRepo.FindNearby(*userLat, *userLon, radiusKm, limit)
		} else {
			// Fallback to search if location not provided
			articles, err = s.articleRepo.SearchByText(intentData.SearchQuery, limit)
		}

	case "category":
		// Category-based filtering
		if intentData.Category != "" {
			articles, err = s.articleRepo.FindByCategory(intentData.Category, limit)
		} else {
			// Fallback to search
			articles, err = s.articleRepo.SearchByText(intentData.SearchQuery, limit)
		}

	case "source":
		// Source-based filtering
		if intentData.Source != "" {
			articles, err = s.articleRepo.FindBySource(intentData.Source, limit)
		} else {
			// Fallback to search
			articles, err = s.articleRepo.SearchByText(intentData.SearchQuery, limit)
		}

	case "score":
		// High-relevance articles
		minScore := intentData.MinScore
		if minScore == 0 {
			minScore = 0.7 // Default threshold
		}
		articles, err = s.articleRepo.FindByMinScore(minScore, limit)

	case "search":
		fallthrough
	default:
		// Full-text search as fallback
		articles, err = s.articleRepo.SearchByText(intentData.SearchQuery, limit)
	}

	if err != nil {
		return nil, intentData, fmt.Errorf("failed to retrieve articles: %w", err)
	}

	return articles, intentData, nil
}
