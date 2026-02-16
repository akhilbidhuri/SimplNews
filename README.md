# SimplNews - Contextual News Data Retrieval System

A Go-based backend system that provides intelligent news article retrieval using LLM-powered intent extraction, geospatial search, and contextual enrichment.

## Features

- 🧠 **LLM-Powered Query Processing** - Natural language queries with GPT-3.5-turbo-16k intent extraction
- 🗺️ **Geospatial Search** - PostGIS-powered location-based article discovery
- 🔍 **Full-Text Search** - PostgreSQL GIN-indexed search with stemming
- 📊 **Trending Analytics** - Location-based trending articles with user engagement simulation
- ⚡ **Fast & Efficient** - In-memory caching, optimized queries, ~22MB Docker image

## Architecture

- **Language**: Go 1.25
- **Database**: PostgreSQL 15 + PostGIS 3.3
- **LLM**: OpenAI GPT-3.5-turbo-16k (intent extraction & summarization)
- **Framework**: Chi Router (lightweight, idiomatic)
- **Caching**: In-memory (sync.Map) + PostgreSQL persistence

## Quick Start

### Prerequisites
- Docker & Docker Compose
- OpenAI API key ([Get one here](https://platform.openai.com/api-keys))

### 1. Clone and Configure

```bash
# Clone the repository
cd SimplNews

# Set your OpenAI API key
export OPENAI_API_KEY=sk-your-actual-key-here

# Or create .env file
cp .env.example .env
# Edit .env and add your OpenAI API key
```

### 2. Start Services

```bash
# Start PostgreSQL and API
docker-compose up -d

# Verify services are running
docker-compose ps
```

### 3. Run Database Migrations

```bash
# Apply migrations
docker exec simplnews-postgres psql -U simplnews_user -d simplnews -f /app/migrations/001_create_articles.up.sql
docker exec simplnews-postgres psql -U simplnews_user -d simplnews -f /app/migrations/002_create_summaries.up.sql
docker exec simplnews-postgres psql -U simplnews_user -d simplnews -f /app/migrations/003_create_events.up.sql
```

### 4. Test API

```bash
# Health check
curl http://localhost:8080/health

# Expected: {"status":"ok","message":"SimplNews API is running"}
```

## API Endpoints

### Main Endpoint (LLM-Powered)

**POST** `/api/v1/news/query` - Natural language query processing

```bash
curl -X POST http://localhost:8080/api/v1/news/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "Latest tech news from Reuters",
    "user_location": null
  }'
```

### Direct Endpoints

- **GET** `/api/v1/news/category?category=world&limit=5` - Filter by category
- **GET** `/api/v1/news/score?min_score=0.7&limit=5` - High relevance articles
- **GET** `/api/v1/news/search?query=election&limit=5` - Full-text search
- **GET** `/api/v1/news/source?source=News18&limit=5` - Filter by source
- **GET** `/api/v1/news/nearby?latitude=17.9&longitude=77.4&radius_km=100&limit=5` - Geospatial search
- **GET** `/api/v1/news/trending?latitude=17.9&longitude=77.4&limit=5` - Trending near location (bonus)

## Development

### Local Development (without Docker)

```bash
# Start only PostgreSQL
docker-compose up -d postgres

# Run API locally
go run cmd/main.go
```

### Build Locally

```bash
go build -o simplnews cmd/main.go
./simplnews
```

### Run Tests

```bash
go test ./...
```

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

## How It Works

### Query Processing Flow

1. **User sends natural language query**: "Tech news from Reuters near Mumbai"
2. **LLM extracts entities & intent**:
   - Entities: `{organizations: ["Reuters"], locations: ["Mumbai"]}`
   - Intent: `source` (Reuters is a news source)
3. **System routes to appropriate endpoint**: Source filter → nearby filter
4. **Results enriched with LLM summaries**: Each article gets a 2-3 sentence summary
5. **Return top 5 articles** in JSON format

### Predefined Entity Labels

- **people** - Names of individuals
- **organizations** - Companies, news sources
- **locations** - Geographic places
- **products** - Specific products

### Intent Mapping

| Intent | Trigger | Endpoint Used |
|--------|---------|---------------|
| `nearby` | Location entity mentioned | Geospatial search |
| `source` | News source organization mentioned | Source filter |
| `category` | Category keywords (tech, politics, sports) | Category filter |
| `score` | Importance words (trending, breaking, top) | High score filter |
| `search` | Default fallback | Full-text search |

## Docker Details

- **Image Size**: ~22.6MB (multi-stage build)
- **Base**: Alpine Linux
- **Health Checks**: Automatic for both PostgreSQL and API
- **Volumes**: Persistent PostgreSQL data

See [DOCKER_SETUP.md](DOCKER_SETUP.md) for detailed Docker instructions.

## Configuration

Configuration via environment variables (.env) and config.yaml:

```env
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=simplnews
OPENAI_API_KEY=sk-your-key
SERVER_PORT=8080
LLM_SUMMARY_MODEL=gpt-3.5-turbo-16k
LLM_INTENT_MODEL=gpt-3.5-turbo-16k
```

## Troubleshooting

### Port already in use
```bash
# PostgreSQL (5432)
lsof -i :5432
kill -9 <PID>

# API (8080)
lsof -i :8080
kill -9 <PID>
```

### Database not initializing
```bash
docker-compose down -v
docker-compose up -d
```

### View logs
```bash
docker-compose logs -f
docker-compose logs simplnews-api
docker-compose logs postgres
```

## Next Steps

1. **Load News Data**: Import 2000 articles from JSON (Phase 2)
2. **Implement Repositories**: Database query layer (Phase 3)
3. **Add LLM Integration**: Intent extraction & summarization (Phase 4)
4. **Build HTTP Endpoints**: All 6 API endpoints (Phase 5-6)
5. **Add Trending System**: User events & trending calculation (Phase 7)

## Tech Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| Backend | Go 1.25 | Fast, concurrent, type-safe |
| Database | PostgreSQL 15 + PostGIS 3.3 | Geospatial queries, full-text search |
| LLM | OpenAI GPT-3.5-turbo-16k | Intent extraction, summarization |
| Router | Chi Router | Lightweight HTTP routing |
| Logger | Zap | Structured JSON logging |
| Config | Viper | Configuration management |

## License

MIT

## Contributing

Phase 1 (Foundation) ✅ Complete
- [x] Project structure
- [x] Configuration & logging
- [x] Database migrations
- [x] Docker setup

Phase 2-8 - In Progress
