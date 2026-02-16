#!/bin/bash

# SimplNews API Endpoint Test Script
# Tests all 6 API endpoints with sample queries

BASE_URL="http://localhost:8080"
BOLD='\033[1m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BOLD}================================${NC}"
echo -e "${BOLD}SimplNews API Endpoint Tests${NC}"
echo -e "${BOLD}================================${NC}\n"

# Test 1: Category Endpoint
echo -e "${GREEN}1. Category Endpoint${NC} - Filter by category"
echo -e "${BLUE}GET /api/v1/news/category?category=technology&limit=3${NC}"
curl -s "${BASE_URL}/api/v1/news/category?category=technology&limit=3" | jq '{
  status,
  count,
  articles: [.data[] | {title, source: .source_name, category}]
}'
echo -e "\n"

# Test 2: Score Endpoint
echo -e "${GREEN}2. Score Endpoint${NC} - Filter by minimum relevance score"
echo -e "${BLUE}GET /api/v1/news/score?min_score=0.85&limit=3${NC}"
curl -s "${BASE_URL}/api/v1/news/score?min_score=0.85&limit=3" | jq '{
  status,
  count,
  articles: [.data[] | {title, score: .relevance_score}]
}'
echo -e "\n"

# Test 3: Search Endpoint (Full-text search)
echo -e "${GREEN}3. Search Endpoint${NC} - PostgreSQL full-text search"
echo -e "${BLUE}GET /api/v1/news/search?query=India&limit=3${NC}"
curl -s "${BASE_URL}/api/v1/news/search?query=India&limit=3" | jq '{
  status,
  count,
  articles: [.data[] | {title, source: .source_name}]
}'
echo -e "\n"

# Test 4: Source Endpoint
echo -e "${GREEN}4. Source Endpoint${NC} - Filter by news source"
echo -e "${BLUE}GET /api/v1/news/source?source=Reuters&limit=3${NC}"
curl -s "${BASE_URL}/api/v1/news/source?source=Reuters&limit=3" | jq '{
  status,
  count,
  articles: [.data[] | {title, source: .source_name, publication_date}]
}'
echo -e "\n"

# Test 5: Nearby Endpoint (Geospatial)
echo -e "${GREEN}5. Nearby Endpoint${NC} - PostGIS geospatial search (Mumbai)"
echo -e "${BLUE}GET /api/v1/news/nearby?latitude=19.0760&longitude=72.8777&radius_km=50&limit=3${NC}"
curl -s "${BASE_URL}/api/v1/news/nearby?latitude=19.0760&longitude=72.8777&radius_km=50&limit=3" | jq '{
  status,
  count,
  articles: [.data[] | {title, location: {lat: .latitude, lon: .longitude}}]
}'
echo -e "\n"

# Test 6: Query Endpoint (LLM-powered)
echo -e "${GREEN}6. Query Endpoint${NC} - LLM-powered intent extraction"
echo -e "${YELLOW}Note: Requires valid OpenAI API key with available credits${NC}"

echo -e "\n${BLUE}Example 1: Source intent${NC}"
echo -e "POST /api/v1/news/query"
echo -e 'Body: {"query": "Technology news from Reuters"}'
curl -s -X POST "${BASE_URL}/api/v1/news/query" \
  -H "Content-Type: application/json" \
  -d '{"query": "Technology news from Reuters"}' | jq '{
  status,
  intent,
  count,
  articles: [.data[]? | {title, source: .source_name}] | .[0:2]
}'
echo -e "\n"

echo -e "${BLUE}Example 2: Category intent${NC}"
echo -e "POST /api/v1/news/query"
echo -e 'Body: {"query": "Show me sports news"}'
curl -s -X POST "${BASE_URL}/api/v1/news/query" \
  -H "Content-Type: application/json" \
  -d '{"query": "Show me sports news"}' | jq '{
  status,
  intent,
  count,
  articles: [.data[]? | {title, category}] | .[0:2]
}'
echo -e "\n"

echo -e "${BLUE}Example 3: Nearby intent${NC}"
echo -e "POST /api/v1/news/query"
echo -e 'Body: {"query": "News near Mumbai", "user_location": {"latitude": 19.0760, "longitude": 72.8777}}'
curl -s -X POST "${BASE_URL}/api/v1/news/query" \
  -H "Content-Type: application/json" \
  -d '{"query": "News near Mumbai", "user_location": {"latitude": 19.0760, "longitude": 72.8777}}' | jq '{
  status,
  intent,
  count,
  articles: [.data[]? | {title, location: {lat: .latitude, lon: .longitude}}] | .[0:2]
}'
echo -e "\n"

echo -e "${BLUE}Example 4: Search fallback${NC}"
echo -e "POST /api/v1/news/query"
echo -e 'Body: {"query": "artificial intelligence"}'
curl -s -X POST "${BASE_URL}/api/v1/news/query" \
  -H "Content-Type: application/json" \
  -d '{"query": "artificial intelligence"}' | jq '{
  status,
  intent,
  count,
  articles: [.data[]? | .title] | .[0:2]
}'
echo -e "\n"

# Health check
echo -e "${GREEN}Health Check${NC}"
echo -e "${BLUE}GET /health${NC}"
curl -s "${BASE_URL}/health" | jq
echo -e "\n"

echo -e "${BOLD}================================${NC}"
echo -e "${BOLD}Test Complete!${NC}"
echo -e "${BOLD}================================${NC}"
