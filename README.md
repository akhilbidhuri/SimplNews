# SimplNews - Contextual News Data Retrieval System

A Go-based backend system that provides intelligent news article retrieval using LLM-powered intent extraction, geospatial search, and full-text search. Built as a POC with 2000 pre-loaded news articles.

## ✨ Features

- 🧠 **LLM-Powered Query Processing** - Natural language queries with GPT-3.5-turbo-16k intent extraction
- 🗺️ **Geospatial Search** - PostGIS-powered location-based article discovery (nearby endpoint)
- 🔍 **Full-Text Search** - PostgreSQL GIN-indexed search on title + description
- 📰 **2000 Pre-loaded Articles** - Real news data with categories, sources, relevance scores
- 🎯 **6 API Endpoints** - Query processing + 5 direct filtering endpoints
- ⚡ **Fast & Efficient** - Optimized queries with indexes, ~22MB Docker image

## 🏗️ Architecture

- **Language**: Go 1.25
- **Database**: PostgreSQL 15 + PostGIS 3.3
- **LLM**: OpenAI GPT-3.5-turbo-16k (intent extraction & summarization)
- **Framework**: Chi Router (lightweight, idiomatic)
- **Caching**: In-memory (sync.Map) + PostgreSQL persistence
- **Logger**: Zap (structured JSON logging)

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- OpenAI API key ([Get one here](https://platform.openai.com/account/api-keys))

### 1. Clone and Configure

```bash
# Clone the repository
cd SimplNews

# Create .env file with your OpenAI API key
cp .env.example .env
# Edit .env and add your OpenAI API key
nano .env
```

### 2. Start Docker Services

```bash
# Start PostgreSQL and API
docker compose up -d

# Verify services are running
docker compose ps

# Check logs
docker compose logs -f api
```

### 3. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Expected: {"status":"healthy"}
```

The API is automatically initialized with:
- ✅ PostgreSQL database with 2000 articles
- ✅ PostGIS geospatial extension
- ✅ All database migrations applied
- ✅ HTTP server running on port 8080

## 📡 API Endpoints

### 1. LLM-Powered Query Processing (POST)

**POST** `/api/v1/news/query` - Natural language intent extraction & routing

```bash
# Example 1: Source intent (finds articles from Reuters)
curl -X POST http://localhost:8080/api/v1/news/query \
  -H "Content-Type: application/json" \
  -d '{"query": "Technology news from Reuters"}'

# Example 2: Category intent (finds sports articles)
curl -X POST http://localhost:8080/api/v1/news/query \
  -H "Content-Type: application/json" \
  -d '{"query": "Show me sports news"}'

# Example 3: Nearby intent (finds articles near Mumbai)
curl -X POST http://localhost:8080/api/v1/news/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "News near Mumbai",
    "user_location": {"latitude": 19.0760, "longitude": 72.8777}
  }'
```

**Response**:
```json
{
  "status": "success",
  "intent": "source",
  "count": 5,
  "data": [
    {
      "id": "uuid",
      "title": "Article title",
      "description": "Article description",
      "source_name": "Reuters",
      "category": ["technology", "world"],
      "relevance_score": 0.85,
      "publication_date": "2025-03-26T10:00:00Z",
      "latitude": 17.9,
      "longitude": 77.4
    }
  ]
}
```

### 2. Direct Filtering Endpoints (GET)

| Endpoint | Purpose | Example |
|----------|---------|---------|
| `/api/v1/news/category` | Filter by category | `?category=technology&limit=5` |
| `/api/v1/news/score` | High relevance articles | `?min_score=0.8&limit=5` |
| `/api/v1/news/search` | Full-text search | `?query=India&limit=5` |
| `/api/v1/news/source` | Filter by news source | `?source=Reuters&limit=5` |
| `/api/v1/news/nearby` | Geospatial search | `?latitude=19.07&longitude=72.87&radius_km=50&limit=5` |

### 3. Test Scripts

Use the provided test scripts for quick testing:

```bash
# Run all endpoints (summary)
./quick-test.sh all

# Test specific endpoint
./quick-test.sh category
./quick-test.sh search
./quick-test.sh nearby

# Comprehensive test with detailed output
./test-endpoints.sh
```

## 💻 Development

### Local Development (without Docker)

```bash
# Start only PostgreSQL
docker compose up -d postgres

# Wait for database to be ready, then run API locally
go run cmd/main.go
```

### Build Locally

```bash
# Build the API binary
go build -o bin/simplnews cmd/main.go

# Build the data loader
go build -o bin/loader cmd/loader/main.go

# Run the API
./bin/simplnews

# Load data (if articles not already loaded)
./bin/loader
```

### Load News Data

The system comes with 2000 pre-loaded articles. To reload data:

```bash
go run cmd/loader/main.go
```

The loader will:
1. Parse `news_data.json` (2000 articles)
2. Validate article data
3. Batch insert into PostgreSQL (100 per transaction)
4. Create PostGIS geometry points for geospatial queries

## Project Structure

```
SimplNews/
├── cmd/
│   ├── main.go              # API server entry point
│   ├── loader/              # Data loading utility
│   └── simulator/           # Event simulator
├── internal/
│   ├── api/                 # HTTP handlers & middleware
│   ├── domain/              # Business logic & services
│   ├── repository/          # Data access layer
│   └── pkg/                 # Shared packages
├── migrations/              # Database schema
├── configs/                 # Configuration files
└── docker-compose.yml       # Docker setup
```

## 🔄 How It Works

### Query Processing Flow

```
User Query
    ↓
LLM Intent Extraction (GPT-3.5-turbo-16k)
    ↓
Entity Detection & Analysis
    ├─ People, Organizations, Locations, Products
    └─ Concepts (abstract ideas)
    ↓
Intent Determination
    ↓
Route to Appropriate Repository Method
    ↓
Return Top N Articles (JSON)
```

### LLM Entities (Predefined Labels)

The LLM extracts only these 4 entity types:

- **people** - Names of individuals (e.g., "Narendra Modi", "Elon Musk")
- **organizations** - Companies, institutions, news sources (e.g., "Reuters", "BBC", "Google")
- **locations** - Geographic places (e.g., "Mumbai", "Silicon Valley", "New York")
- **products** - Specific products (e.g., "iPhone 15", "ChatGPT")

### Intent Routing

The system automatically determines intent based on detected entities:

| Intent | Trigger | Repository Method | Example Query |
|--------|---------|------------------|---|
| `nearby` | Location entity + user coordinates | `FindNearby()` | "News near Mumbai" |
| `source` | Organization entity (news source) | `FindBySource()` | "Reuters articles about tech" |
| `category` | Category keywords (tech, sports, etc.) | `FindByCategory()` | "Show me technology news" |
| `score` | Importance words (trending, breaking, top) | `FindByMinScore()` | "What's trending?" |
| `search` | Default fallback (no specific intent) | `SearchByText()` | "artificial intelligence" |

### Example: Query Processing

**Input**: "Technology news from Reuters"

1. **LLM Analysis**:
   - Entities: `{organizations: ["Reuters"]}`
   - Concepts: `["technology"]`
   - Intent: `source` (Reuters is a news organization)

2. **System Routes**: `FindBySource("Reuters", limit)`

3. **Database Query**:
   ```sql
   SELECT * FROM articles WHERE source_name = 'Reuters' LIMIT 5
   ```

4. **Response**: Top 5 articles from Reuters

## 🐳 Docker Details

- **Image Size**: ~22.6MB (multi-stage build with Alpine)
- **Base Image**: `golang:1.25-alpine` (builder) → `alpine:latest` (runtime)
- **Health Checks**: Automatic for both PostgreSQL and API
- **Volumes**:
  - `pgdata` - Persistent PostgreSQL data
  - `./configs` - Configuration files
  - `./migrations` - Database migrations
- **Startup**: Automatic migrations applied, 2000 articles pre-loaded

### Docker Compose Commands

```bash
# Start services
docker compose up -d

# View logs
docker compose logs -f
docker compose logs api
docker compose logs postgres

# Stop services
docker compose down

# Stop and remove volumes (fresh start)
docker compose down -v

# Check service status
docker compose ps
```

## ⚙️ Configuration

Configuration via environment variables in `.env`:

```env
# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=simplnews
DATABASE_USER=simplnews_user
DATABASE_PASSWORD=changeme123

# OpenAI API
OPENAI_API_KEY=sk-your-actual-key-here

# Server
SERVER_PORT=8080

# LLM Models
LLM_SUMMARY_MODEL=gpt-3.5-turbo-16k
LLM_INTENT_MODEL=gpt-3.5-turbo-16k
LLM_INTENT_MAX_TOKENS=300
LLM_INTENT_TEMPERATURE=0.1

# Logging
LOG_LEVEL=info
```

### For Docker (.env.docker)

Copy `.env.docker` for use with Docker Compose:

```bash
cp .env.docker .env
# Edit and add your OpenAI API key
```

## 🛠️ Troubleshooting

### OpenAI API Key Issues

**Error**: `OPENAI_API_KEY environment variable is required`

**Solution**: Ensure your OpenAI API key is set in `.env`:
```bash
OPENAI_API_KEY=sk-proj-your-actual-key-here
```

Then restart:
```bash
docker compose down && docker compose up -d
```

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker compose ps postgres

# View database logs
docker compose logs postgres

# Fresh database (removes all data)
docker compose down -v
docker compose up -d
```

### Port Already in Use

```bash
# PostgreSQL (5432)
lsof -i :5432
kill -9 <PID>

# API (8080)
lsof -i :8080
kill -9 <PID>

# Or use different ports in docker-compose.yml
```

### View Full Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f api
docker compose logs -f postgres

# Last 50 lines
docker compose logs --tail 50
```

## 📦 Tech Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Backend** | Go 1.25 | Type-safe, concurrent, fast HTTP server |
| **Database** | PostgreSQL 15 + PostGIS 3.3 | Geospatial queries, full-text search, ACID transactions |
| **LLM** | OpenAI GPT-3.5-turbo-16k | Intent extraction, entity detection, summarization |
| **HTTP Router** | Chi Router v5 | Lightweight, idiomatic Go routing |
| **Logger** | Zap | Structured JSON logging with levels |
| **Config** | Viper + godotenv | Environment-based configuration |
| **Database Driver** | sqlx + pq | PostgreSQL driver with extended functionality |
| **Geospatial** | PostGIS 3.3 | Geographic point/distance calculations |
| **Text Search** | PostgreSQL GIN | Full-text search with English stemming |
| **Container** | Docker & Alpine | Lightweight deployments (~22MB) |

## 📊 Database Schema

### articles
- `id` (UUID) - Primary key
- `title` (text) - Article headline
- `description` (text) - Article body
- `url` (text) - Source URL
- `publication_date` (timestamp) - When published
- `source_name` (varchar) - News source (Reuters, BBC, etc.)
- `category` (text[]) - Array of categories (technology, sports, etc.)
- `relevance_score` (double) - 0.0 to 1.0 relevance ranking
- `location` (geography) - PostGIS point for geospatial queries
- `latitude`/`longitude` (double) - Decimal coordinates
- `created_at` (timestamp) - Record creation time
- `updated_at` (timestamp) - Last update time

**Indexes**:
- `pk_articles` - UUID primary key
- `idx_articles_category` - GIN on category array
- `idx_articles_search` - GIN full-text search (title + description)
- `idx_articles_source_name` - B-tree on source_name
- `idx_articles_relevance_score` - B-tree on relevance_score (DESC)
- `idx_articles_publication_date` - B-tree on publication_date (DESC)
- `idx_articles_location` - GiST on PostGIS geography

## 📝 Sample Queries

```bash
# Category filter - all technology articles
curl "http://localhost:8080/api/v1/news/category?category=technology&limit=10"

# High relevance - trending articles
curl "http://localhost:8080/api/v1/news/score?min_score=0.85&limit=5"

# Full-text search - search for "election"
curl "http://localhost:8080/api/v1/news/search?query=election&limit=10"

# Source filter - articles from Reuters
curl "http://localhost:8080/api/v1/news/source?source=Reuters&limit=5"

# Geospatial - articles within 100km of Mumbai
curl "http://localhost:8080/api/v1/news/nearby?latitude=19.0760&longitude=72.8777&radius_km=100&limit=10"

# LLM-powered - intelligent query processing
curl -X POST http://localhost:8080/api/v1/news/query \
  -H "Content-Type: application/json" \
  -d '{"query": "Tech news near Mumbai"}'
```

## 📚 Project Structure

```
SimplNews/
├── cmd/
│   ├── main.go                  # API server entry point
│   └── loader/
│       └── main.go              # Data loader utility
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP request handlers (6 endpoints)
│   │   ├── routes.go            # Route definitions
│   │   └── server.go            # HTTP server setup
│   ├── domain/
│   │   ├── models/
│   │   │   └── article.go       # Article data model
│   │   └── services/
│   │       ├── llm_service.go   # GPT-3.5 integration
│   │       └── query_service.go # Query processing orchestration
│   ├── repository/
│   │   └── postgres/
│   │       ├── db.go            # Database connection
│   │       └── article_repo.go  # Article queries (5 methods)
│   ├── loader/
│   │   ├── json_loader.go       # Parse news_data.json
│   │   └── db_importer.go       # Batch insert to PostgreSQL
│   └── pkg/
│       ├── config/              # Configuration loading
│       └── logger/              # Structured logging
├── migrations/                  # PostgreSQL schema
├── configs/
│   └── config.yaml              # Default configuration
├── docker-compose.yml           # Docker setup
├── Dockerfile                   # Multi-stage build
├── test-endpoints.sh            # Comprehensive test script
├── quick-test.sh                # Quick endpoint testing
└── README.md                    # This file
```

## 🚧 Implementation Status

### ✅ Completed
- [x] Project structure & configuration
- [x] PostgreSQL database with PostGIS
- [x] Database schema with indexes
- [x] Data loader (2000 articles)
- [x] Article repository (5 query methods)
- [x] LLM service (intent extraction & summarization)
- [x] Query service (intelligent routing)
- [x] 6 HTTP endpoints
- [x] Docker setup with health checks
- [x] Test scripts (comprehensive & quick)

### 🔄 Future Enhancements
- [ ] User event tracking & trending analytics
- [ ] Article caching & optimization
- [ ] Summary caching with TTL
- [ ] Advanced filtering (date ranges, multiple categories)
- [ ] Pagination for large result sets
- [ ] API authentication & rate limiting
- [ ] GraphQL endpoint
- [ ] WebSocket for real-time updates

## 📄 License

MIT

---

**Last Updated**: February 2026
**Version**: 1.0.0 - POC
