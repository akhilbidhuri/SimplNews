package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/akhilbidhuri/SimplNews/internal/domain/models"
	"github.com/akhilbidhuri/SimplNews/internal/loader"
	"github.com/akhilbidhuri/SimplNews/internal/pkg/config"
	"github.com/akhilbidhuri/SimplNews/internal/pkg/logger"
	"github.com/akhilbidhuri/SimplNews/internal/repository/postgres"
)

func main() {
	// Parse command-line flags
	filePath := flag.String("file", "news_data.json", "Path to news_data.json file")
	validate := flag.Bool("validate", true, "Validate articles before importing")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(cfg.Logging.Level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Infow("SimplNews Data Loader starting",
		"file", *filePath,
		"validate", *validate,
	)

	// Connect to database
	db, err := postgres.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalw("Failed to connect to database", "error", err)
	}
	defer db.Close()

	log.Infow("Connected to PostgreSQL", "database", cfg.Database.Name)

	// Load articles from JSON
	startLoad := time.Now()
	log.Infow("Loading articles from JSON", "file", *filePath)
	
	articles, err := loader.LoadArticlesFromJSON(*filePath)
	if err != nil {
		log.Fatalw("Failed to load articles from JSON", "error", err)
	}
	
	loadDuration := time.Since(startLoad)
	log.Infow("Articles loaded from JSON",
		"count", len(articles),
		"duration", loadDuration,
	)

	// Validate articles if requested
	if *validate {
		log.Infow("Validating articles")
		validCount := 0
		invalidCount := 0
		
		validArticles := make([]models.Article, 0, len(articles))
		for i, article := range articles {
			if err := loader.ValidateArticle(&article); err != nil {
				log.Warnw("Invalid article",
					"index", i,
					"id", article.ID,
					"error", err,
				)
				invalidCount++
			} else {
				validArticles = append(validArticles, article)
				validCount++
			}
		}
		
		log.Infow("Validation complete",
			"valid", validCount,
			"invalid", invalidCount,
		)
		
		articles = validArticles
	}

	// Import articles into database
	startImport := time.Now()
	log.Infow("Importing articles to PostgreSQL", "count", len(articles))
	
	if err := loader.ImportArticles(db, articles); err != nil {
		log.Fatalw("Failed to import articles", "error", err)
	}
	
	importDuration := time.Since(startImport)
	log.Infow("Articles imported successfully", "duration", importDuration)

	// Verify import
	count, err := loader.GetArticleCount(db)
	if err != nil {
		log.Errorw("Failed to verify import", "error", err)
	} else {
		log.Infow("Total articles in database", "count", count)
	}

	totalDuration := time.Since(startLoad)
	log.Infow("Data loading complete",
		"total_duration", totalDuration,
		"articles_processed", len(articles),
	)
}
