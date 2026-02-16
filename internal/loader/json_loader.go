package loader

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/akhilbidhuri/SimplNews/internal/domain/models"
)

// LoadArticlesFromJSON reads and parses the news_data.json file
func LoadArticlesFromJSON(filePath string) ([]models.Article, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var articles []models.Article
	decoder := json.NewDecoder(file)
	
	if err := decoder.Decode(&articles); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return articles, nil
}

// ValidateArticle performs basic validation on an article
func ValidateArticle(article *models.Article) error {
	if article.ID == "" {
		return fmt.Errorf("article ID is empty")
	}
	if article.Title == "" {
		return fmt.Errorf("article title is empty")
	}
	if article.SourceName == "" {
		return fmt.Errorf("article source_name is empty")
	}
	if article.RelevanceScore < 0 || article.RelevanceScore > 1 {
		return fmt.Errorf("relevance_score must be between 0 and 1, got: %f", article.RelevanceScore)
	}
	if article.Latitude < -90 || article.Latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90, got: %f", article.Latitude)
	}
	if article.Longitude < -180 || article.Longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180, got: %f", article.Longitude)
	}
	if len(article.Category) == 0 {
		return fmt.Errorf("article must have at least one category")
	}
	return nil
}
