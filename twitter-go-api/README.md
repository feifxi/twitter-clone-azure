# Twitter Clone (Go API)

This is the high-performance, production-ready backend API for the Twitter/X clone, built directly in Go.

It was designed to replace the original Java/Spring Boot backend to provide significantly faster execution, lower memory footprint, and highly optimized database querying using modern Go patterns.

## 🚀 Tech Stack

- **Language:** Go 1.22+
- **Framework:** Gin Web Framework
- **Database:** PostgreSQL
- **Query Builder:** sqlc (Type-safe SQL code generation)
- **Caching:** Redis v9
- **Logging:** Zerolog (Structured JSON Telemetry)
- **Authentication:** Custom JWT authentication (Tokens)
- **Storage:** Azure Blob Storage (for media uploads)
- **Architecture:** Clean Architecture (Handler -> Usecase -> Repository)

## ✨ Key Features & Optimizations

- **Memory Batching (DataLoader Pattern):** Eliminates the N+1 query problem commonly found in social media feeds. The API batches database IDs (like `UserID`, `ParentID`, `RetweetID`) and fetches entire graph relationships in just 2-3 optimal SQL queries.
- **Redis Caching:** Drastically reduces DB load by caching Anonymous Feeds and Trending Hashtags in Redis, while dynamically injecting user-specific `IsLiked` parameters securely from memory.
- **Enterprise Hardening:** Utilizes Zerolog for structured JSON observability, enforces `Read/WriteTimeouts` to mitigate Slowloris DDoS attacks, and cleanly masks internal database constraints into proper API Conflict responses.
- **Efficient JSON Serialization:** Direct mapping of nested structs with exact camelCase JSON tags corresponding to the frontend expectations.
- **Full Text Search:** Leverages PostgreSQL's native `to_tsvector` and `ts_rank` for blazing fast full-text tweet searching and hashtag prefix matching.
- **Media Uploads:** Multipart-form parsing directly streaming to Azure Blob storage.

## 🛠️ Getting Started (Local Development)

### Prerequisites

- Go 1.22 or higher
- PostgreSQL running locally (can be started via Docker Compose in the project root)
- [sqlc](https://sqlc.dev/) (for database model generation)
- [golang-migrate](https://github.com/golang-migrate/migrate) (for database migrations)

### 1. Database Setup & Migrations

Ensure your PostgreSQL database `twitter_db` is running.
```bash
# Run database migrations
make migrateup
```

*Note: The database configuration expects `postgres://root:rootpass@localhost:5432/twitter_db?sslmode=disable` by default.*

### 2. Configuration

Create a `app.env` file in the root of `twitter-go-api` based on required environment variables:
```env
DB_SOURCE=postgres://root:rootpass@localhost:5432/twitter_db?sslmode=disable
SERVER_ADDRESS=0.0.0.0:8080
TOKEN_SYMMETRIC_KEY=12345678901234567890123456789012
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=24h
AZURE_STORAGE_ACCOUNT_NAME=your_account_name
AZURE_STORAGE_ACCOUNT_KEY=your_access_key
AZURE_CONTAINER_NAME=twitter-media
REDIS_ADDRESS=redis://localhost:6379
REDIS_PASSWORD=
```

### 3. Running the API

```bash
# Generate SQLC models (if you modify queries in db/query/*.sql)
make sqlc

# Run the API
make run
```

The API will be available at `http://localhost:8080`.

## 📂 Project Structure

- `/cmd/api` - The main application entry point.
- `/db` - Contains `migration` files and `query` SQL schemas. `sqlc` reads these to generate the Go database repository.
- `/internal/db` - The generated `sqlc` database interaction code. Do not edit manually.
- `/internal/server` - Gin HTTP handlers, routes, request parsing, and JSON response mapping.
- `/internal/usecase` - The core business logic, including the DataLoader batching mechanisms.
- `/internal/config` - Viper configuration management.
- `/internal/token` - JWT token creation and verification.
- `/internal/service` - Third-party integrations (like Azure Blob Storage).
