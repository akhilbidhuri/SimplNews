package postgres

import (
	"database/sql"
	"fmt"

	"github.com/akhilbidhuri/SimplNews/internal/domain/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ArticleRepository struct {
	db *sqlx.DB
}

func NewArticleRepository(db *sqlx.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// FindByCategory returns articles matching a specific category
func (r *ArticleRepository) FindByCategory(category string, limit int) ([]models.Article, error) {
	query := `
		SELECT id, title, description, url, publication_date, source_name,
		       category, relevance_score, latitude, longitude, created_at, updated_at
		FROM articles
		WHERE $1 = ANY(category)
		ORDER BY publication_date DESC, relevance_score DESC
		LIMIT $2
	`

	var articles []models.Article
	err := r.db.Select(&articles, query, category, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find articles by category: %w", err)
	}

	return articles, nil
}

// FindByMinScore returns articles with relevance score >= minScore
func (r *ArticleRepository) FindByMinScore(minScore float64, limit int) ([]models.Article, error) {
	query := `
		SELECT id, title, description, url, publication_date, source_name,
		       category, relevance_score, latitude, longitude, created_at, updated_at
		FROM articles
		WHERE relevance_score >= $1
		ORDER BY relevance_score DESC, publication_date DESC
		LIMIT $2
	`

	var articles []models.Article
	err := r.db.Select(&articles, query, minScore, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find articles by score: %w", err)
	}

	return articles, nil
}

// SearchByText performs full-text search on title and description
func (r *ArticleRepository) SearchByText(searchQuery string, limit int) ([]models.Article, error) {
	query := `
		SELECT id, title, description, url, publication_date, source_name,
		       category, relevance_score, latitude, longitude, created_at, updated_at
		FROM articles
		WHERE to_tsvector('english', title || ' ' || description)
		      @@ plainto_tsquery('english', $1)
		ORDER BY ts_rank(
		           to_tsvector('english', title || ' ' || description),
		           plainto_tsquery('english', $1)
		       ) DESC, relevance_score DESC
		LIMIT $2
	`

	var articles []models.Article
	err := r.db.Select(&articles, query, searchQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search articles: %w", err)
	}

	return articles, nil
}

// FindBySource returns articles from a specific news source
func (r *ArticleRepository) FindBySource(source string, limit int) ([]models.Article, error) {
	query := `
		SELECT id, title, description, url, publication_date, source_name,
		       category, relevance_score, latitude, longitude, created_at, updated_at
		FROM articles
		WHERE source_name = $1
		ORDER BY publication_date DESC, relevance_score DESC
		LIMIT $2
	`

	var articles []models.Article
	err := r.db.Select(&articles, query, source, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find articles by source: %w", err)
	}

	return articles, nil
}

// FindNearby returns articles within radiusKm of the given coordinates
func (r *ArticleRepository) FindNearby(lat, lon, radiusKm float64, limit int) ([]models.Article, error) {
	query := `
		SELECT id, title, description, url, publication_date, source_name,
		       category, relevance_score, latitude, longitude, created_at, updated_at
		FROM articles
		WHERE ST_DWithin(
		    location,
		    ST_MakePoint($2, $1)::geography,
		    $3 * 1000
		)
		ORDER BY ST_Distance(
		           location,
		           ST_MakePoint($2, $1)::geography
		       ) ASC
		LIMIT $4
	`

	var articles []models.Article
	err := r.db.Select(&articles, query, lat, lon, radiusKm, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearby articles: %w", err)
	}

	return articles, nil
}

// GetAllCategories returns unique categories present in the database
func (r *ArticleRepository) GetAllCategories() ([]string, error) {
	query := `
		SELECT DISTINCT unnest(category) as category
		FROM articles
		ORDER BY category
	`

	var categories []string
	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, nil
}

// GetAllSources returns unique news sources present in the database
func (r *ArticleRepository) GetAllSources() ([]string, error) {
	query := `
		SELECT DISTINCT source_name
		FROM articles
		ORDER BY source_name
	`

	var sources []string
	err := r.db.Select(&sources, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sources: %w", err)
	}

	return sources, nil
}

// GetByID returns a single article by ID
func (r *ArticleRepository) GetByID(id string) (*models.Article, error) {
	query := `
		SELECT id, title, description, url, publication_date, source_name,
		       category, relevance_score, latitude, longitude, created_at, updated_at
		FROM articles
		WHERE id = $1
	`

	var article models.Article
	err := r.db.Get(&article, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("article not found")
		}
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	return &article, nil
}

// GetByIDs returns multiple articles by their IDs
func (r *ArticleRepository) GetByIDs(ids []string) ([]models.Article, error) {
	query := `
		SELECT id, title, description, url, publication_date, source_name,
		       category, relevance_score, latitude, longitude, created_at, updated_at
		FROM articles
		WHERE id = ANY($1)
	`

	var articles []models.Article
	err := r.db.Select(&articles, query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("failed to get articles by IDs: %w", err)
	}

	return articles, nil
}
