# SimplNews Docker Setup Guide

## Prerequisites
- Docker and Docker Compose installed
- OpenAI API key

## Quick Start

### 1. Set OpenAI API Key
Before starting the services, update the OpenAI API key in docker-compose.yml or create a .env file:

```bash
# Option A: Create .env file
cp .env.docker .env
# Then edit .env and replace the OPENAI_API_KEY value

# Option B: Export environment variable
export OPENAI_API_KEY=sk-your-actual-key-here
```

### 2. Start Services
```bash
# Build and start PostgreSQL and API services
docker-compose up -d

# Check if services are running
docker-compose ps

# View logs
docker-compose logs -f
```

### 3. Verify Database is Ready
```bash
# Wait for PostgreSQL to be healthy
docker-compose ps postgres
# Status should show "healthy" before proceeding

# Run migrations
docker exec simplnews-postgres psql -U simplnews_user -d simplnews -f /migrations/001_create_articles.up.sql
```

### 4. Test API Health
```bash
# Check API is running
curl http://localhost:8080/health

# Expected response:
# {"status":"ok","message":"SimplNews API is running"}
```

## Service Details

### PostgreSQL (simplnews-postgres)
- **Host**: localhost:5432 (from host), postgres:5432 (from Docker)
- **Database**: simplnews
- **User**: simplnews_user
- **Password**: changeme123
- **Image**: postgis/postgis:15-3.3

### API Service (simplnews-api)
- **Host**: localhost:8080
- **Health Check**: GET /health
- **Dependencies**: Waits for PostgreSQL to be healthy before starting

## Common Commands

```bash
# View all logs
docker-compose logs

# View API logs only
docker-compose logs simplnews-api

# View PostgreSQL logs only
docker-compose logs postgres

# Stop services
docker-compose down

# Stop and remove volumes (WARNING: deletes database)
docker-compose down -v

# Rebuild API image after code changes
docker-compose build --no-cache api

# Execute command in running container
docker exec simplnews-postgres psql -U simplnews_user -d simplnews -c "SELECT COUNT(*) FROM articles;"

# Open PostgreSQL shell
docker exec -it simplnews-postgres psql -U simplnews_user -d simplnews

# Open API container shell
docker exec -it simplnews-api sh
```

## Troubleshooting

### Port 5432 already in use
```bash
# Find and stop the service using port 5432
lsof -i :5432
kill -9 <PID>

# Or change PostgreSQL port in docker-compose.yml
# Change: ports: - "5433:5432"
```

### Port 8080 already in use
```bash
# Find and stop the service using port 8080
lsof -i :8080
kill -9 <PID>

# Or change API port in docker-compose.yml
# Change: ports: - "8081:8080"
```

### Database not initializing
```bash
# Remove volumes and restart
docker-compose down -v
docker-compose up -d

# Wait for healthy status
docker-compose ps
```

### API not starting
```bash
# Check logs
docker-compose logs simplnews-api

# Common issues:
# 1. OPENAI_API_KEY not set
# 2. PostgreSQL not healthy yet
# 3. Invalid .env configuration
```

## Development Workflow

### Local Development (without Docker)
```bash
# Start only PostgreSQL
docker-compose up -d postgres

# Run API locally
go run cmd/main.go
```

### Production-like Testing
```bash
# Start everything with Docker
docker-compose up -d

# Test all endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/news/category?category=world
```

## Environment Variables

All environment variables are defined in `docker-compose.yml` or can be overridden with `.env` file:

```env
DATABASE_HOST=postgres
DATABASE_PORT=5432
DATABASE_NAME=simplnews
DATABASE_USER=simplnews_user
DATABASE_PASSWORD=changeme123
OPENAI_API_KEY=sk-your-key
SERVER_PORT=8080
LOG_LEVEL=info
```

## Image Details

The Docker image is built using multi-stage build:
- **Size**: ~22.6MB (very lightweight)
- **Base Image**: alpine:latest
- **Includes**: Go binary, PostgreSQL client, CA certificates

## Next Steps

After Docker setup is verified:
1. Load news data: `go run cmd/loader/main.go`
2. Run tests: `docker-compose exec api go test ./...`
3. Deploy to production: Push image to registry
