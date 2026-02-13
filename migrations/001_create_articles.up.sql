-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create articles table
CREATE TABLE articles (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    url TEXT NOT NULL,
    publication_date TIMESTAMPTZ NOT NULL,
    source_name VARCHAR(255) NOT NULL,
    category TEXT[] NOT NULL,
    relevance_score FLOAT NOT NULL CHECK (relevance_score >= 0 AND relevance_score <= 1),
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_articles_category ON articles USING GIN(category);
CREATE INDEX idx_articles_source_name ON articles(source_name);
CREATE INDEX idx_articles_relevance_score ON articles(relevance_score DESC);
CREATE INDEX idx_articles_publication_date ON articles(publication_date DESC);
CREATE INDEX idx_articles_location ON articles USING GIST(location);
CREATE INDEX idx_articles_search ON articles USING GIN(to_tsvector('english', title || ' ' || description));
