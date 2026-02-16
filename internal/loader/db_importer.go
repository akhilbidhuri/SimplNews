package loader

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/akhilbidhuri/SimplNews/internal/domain/models"
)

const batchSize = 100

// ImportArticles imports articles into PostgreSQL in batches
func ImportArticles(db *sqlx.DB, articles []models.Article) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare the insert statement
	query := `
		INSERT INTO articles (
			id, title, description, url, publication_date,
			source_name, category, relevance_score,
			location, latitude, longitude
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			ST_MakePoint($10, $9)::geography, $9, $10
		)
		ON CONFLICT (id) DO NOTHING
	`

	stmt, err := tx.Preparex(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert articles in batches
	inserted := 0
	skipped := 0

	for i, article := range articles {
		result, err := stmt.Exec(
			article.ID,
			article.Title,
			article.Description,
			article.URL,
			article.PublicationDate,
			article.SourceName,
			pq.Array(article.Category), // PostgreSQL array type
			article.RelevanceScore,
			article.Latitude,
			article.Longitude,
		)
		
		if err != nil {
			return fmt.Errorf("failed to insert article %d (ID: %s): %w", i, article.ID, err)
		}

		rows, _ := result.RowsAffected()
		if rows > 0 {
			inserted++
		} else {
			skipped++
		}

		// Commit in batches for better performance
		if (i+1)%batchSize == 0 {
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("failed to commit batch at index %d: %w", i, err)
			}
			
			// Start new transaction for next batch
			tx, err = db.Beginx()
			if err != nil {
				return fmt.Errorf("failed to begin new transaction: %w", err)
			}
			
			stmt, err = tx.Preparex(query)
			if err != nil {
				return fmt.Errorf("failed to prepare statement for new batch: %w", err)
			}
		}
	}

	// Commit remaining articles
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit final batch: %w", err)
	}

	fmt.Printf("Import complete: %d inserted, %d skipped (duplicates)\n", inserted, skipped)
	return nil
}

// GetArticleCount returns the count of articles in the database
func GetArticleCount(db *sqlx.DB) (int, error) {
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM articles")
	if err != nil {
		return 0, fmt.Errorf("failed to get article count: %w", err)
	}
	return count, nil
}
