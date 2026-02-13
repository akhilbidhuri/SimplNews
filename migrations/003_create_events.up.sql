-- Create table for simulated user events (trending system)
CREATE TABLE user_events (
    id BIGSERIAL PRIMARY KEY,
    article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL,
    user_location GEOGRAPHY(POINT, 4326),
    latitude FLOAT,
    longitude FLOAT,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_user_events_article_id ON user_events(article_id);
CREATE INDEX idx_user_events_timestamp ON user_events(timestamp DESC);
CREATE INDEX idx_user_events_location ON user_events USING GIST(user_location);
