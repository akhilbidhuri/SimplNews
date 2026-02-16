#!/bin/bash

# Quick test script for SimplNews API
# Usage: ./quick-test.sh [endpoint]
# Available endpoints: category, score, search, source, nearby, query, all

BASE_URL="http://localhost:8080"

case "$1" in
  category)
    echo "Testing Category endpoint..."
    curl -s "${BASE_URL}/api/v1/news/category?category=technology&limit=3" | jq
    ;;

  score)
    echo "Testing Score endpoint..."
    curl -s "${BASE_URL}/api/v1/news/score?min_score=0.8&limit=3" | jq
    ;;

  search)
    echo "Testing Search endpoint..."
    curl -s "${BASE_URL}/api/v1/news/search?query=India&limit=3" | jq
    ;;

  source)
    echo "Testing Source endpoint..."
    curl -s "${BASE_URL}/api/v1/news/source?source=Reuters&limit=3" | jq
    ;;

  nearby)
    echo "Testing Nearby endpoint (Mumbai)..."
    curl -s "${BASE_URL}/api/v1/news/nearby?latitude=19.0760&longitude=72.8777&radius_km=50&limit=3" | jq
    ;;

  query)
    echo "Testing Query endpoint (LLM-powered)..."
    echo "Query: 'Technology news from Reuters'"
    curl -s -X POST "${BASE_URL}/api/v1/news/query" \
      -H "Content-Type: application/json" \
      -d '{"query": "Technology news from Reuters"}' | jq
    ;;

  health)
    echo "Testing Health endpoint..."
    curl -s "${BASE_URL}/health" | jq
    ;;

  all)
    echo "=== Testing All Endpoints ==="
    echo ""
    echo "1. Category:"
    curl -s "${BASE_URL}/api/v1/news/category?category=technology&limit=2" | jq -c '{status, count}'
    echo ""
    echo "2. Score:"
    curl -s "${BASE_URL}/api/v1/news/score?min_score=0.8&limit=2" | jq -c '{status, count}'
    echo ""
    echo "3. Search:"
    curl -s "${BASE_URL}/api/v1/news/search?query=India&limit=2" | jq -c '{status, count}'
    echo ""
    echo "4. Source:"
    curl -s "${BASE_URL}/api/v1/news/source?source=Reuters&limit=2" | jq -c '{status, count}'
    echo ""
    echo "5. Nearby:"
    curl -s "${BASE_URL}/api/v1/news/nearby?latitude=19.0760&longitude=72.8777&radius_km=50&limit=2" | jq -c '{status, count}'
    echo ""
    echo "6. Query (LLM):"
    curl -s -X POST "${BASE_URL}/api/v1/news/query" \
      -H "Content-Type: application/json" \
      -d '{"query": "Technology news"}' | jq -c '{status, intent, count}'
    ;;

  *)
    echo "SimplNews API Quick Test"
    echo ""
    echo "Usage: $0 [endpoint]"
    echo ""
    echo "Available endpoints:"
    echo "  category  - Filter by category"
    echo "  score     - Filter by relevance score"
    echo "  search    - Full-text search"
    echo "  source    - Filter by news source"
    echo "  nearby    - Geospatial search"
    echo "  query     - LLM-powered query (requires OpenAI credits)"
    echo "  health    - Health check"
    echo "  all       - Test all endpoints (summary)"
    echo ""
    echo "Examples:"
    echo "  $0 category"
    echo "  $0 search"
    echo "  $0 all"
    ;;
esac
