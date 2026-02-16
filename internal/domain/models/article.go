package models

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// Article represents a news article with geospatial data
type Article struct {
	ID              string         `json:"id" db:"id"`
	Title           string         `json:"title" db:"title"`
	Description     string         `json:"description" db:"description"`
	URL             string         `json:"url" db:"url"`
	PublicationDate time.Time      `json:"publication_date" db:"publication_date"`
	SourceName      string         `json:"source_name" db:"source_name"`
	Category        pq.StringArray `json:"category" db:"category"`
	RelevanceScore  float64        `json:"relevance_score" db:"relevance_score"`
	Latitude        float64        `json:"latitude" db:"latitude"`
	Longitude       float64        `json:"longitude" db:"longitude"`
	CreatedAt       time.Time      `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at,omitempty" db:"updated_at"`
}

// ArticleWithSummary extends Article with LLM-generated summary
type ArticleWithSummary struct {
	Article
	LLMSummary string `json:"llm_summary,omitempty"`
}

// Location represents a geographic coordinate
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// CustomTime handles the custom timestamp format from news_data.json
type CustomTime struct {
	time.Time
}

// UnmarshalJSON parses timestamps in the format "2025-03-26T04:46:55"
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	// Remove quotes
	s = s[1 : len(s)-1]

	// Try parsing with timezone first
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		// If that fails, try without timezone (append Z)
		t, err = time.Parse(time.RFC3339, s+"Z")
		if err != nil {
			// Try custom format without timezone
			t, err = time.Parse("2006-01-02T15:04:05", s)
			if err != nil {
				return err
			}
		}
	}

	ct.Time = t
	return nil
}

// ArticleJSON is used for unmarshaling from JSON with custom time format
type ArticleJSON struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	URL             string     `json:"url"`
	PublicationDate CustomTime `json:"publication_date"`
	SourceName      string     `json:"source_name"`
	Category        []string   `json:"category"`
	RelevanceScore  float64    `json:"relevance_score"`
	Latitude        float64    `json:"latitude"`
	Longitude       float64    `json:"longitude"`
}

// UnmarshalJSON implements custom JSON unmarshaling for Article
func (a *Article) UnmarshalJSON(data []byte) error {
	var aj ArticleJSON
	if err := json.Unmarshal(data, &aj); err != nil {
		return err
	}

	a.ID = aj.ID
	a.Title = aj.Title
	a.Description = aj.Description
	a.URL = aj.URL
	a.PublicationDate = aj.PublicationDate.Time
	a.SourceName = aj.SourceName
	a.Category = aj.Category
	a.RelevanceScore = aj.RelevanceScore
	a.Latitude = aj.Latitude
	a.Longitude = aj.Longitude

	return nil
}
