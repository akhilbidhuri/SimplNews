-- Create table for LLM-generated summaries
CREATE TABLE article_summaries (
    article_id UUID PRIMARY KEY REFERENCES articles(id) ON DELETE CASCADE,
    summary TEXT NOT NULL,
    generated_at TIMESTAMPTZ DEFAULT NOW(),
    llm_model VARCHAR(100) NOT NULL
);

CREATE INDEX idx_article_summaries_generated_at ON article_summaries(generated_at DESC);
