# Ready-Go CLI

A CLI tool to scaffold production-ready Go projects with Fiber v3, SQLC, and a practical architecture.

## Features

- **Fiber v3**: Latest version with built-in timeout support and improved performance
- **Type-Safe SQL**: SQLC generates Go code from your SQL queries
- **Clean Architecture**: Handlers → SQLC Models → Database (no unnecessary layers)
- **Docker Ready**: MySQL, Redis, and Kafka pre-configured
- **Request Logging**: Built-in HTTP logging with structured output
- **Error Handling**: Structured error responses with error codes
- **Domain Organization**: Handlers organized by domain in subdirectories

## Installation

```bash
# Install from source
go install github.com/muazwzxv/ready-go-cli/cmd/ready-go@latest

# Or clone and build
git clone https://github.com/muazwzxv/ready-go-cli.git
cd ready-go-cli
go build -o ready-go ./cmd/ready-go
mv ready-go ~/go/bin/
```

## Quick Start

```bash
# Create a new project
ready-go new my-api --module github.com/mycompany/my-api

cd my-api

# Start infrastructure
make docker-up

# Run migrations
make migrate-up

# Generate SQLC models
make sqlc-generate

# Start the API
make run-api
```

Visit `http://localhost:8080`

## Usage

```bash
ready-go new <project-name> [flags]

Flags:
  --module, -m    Go module path (default: github.com/username/<project>)
  --port          Server port (default: 8080)
  --db-port       MySQL port (default: 3306)
  --redis-port    Redis port (default: 6379)
  --kafka-port    Kafka port (default: 9092)
  --sample-name   Sample entity name (default: User)
```

## Generated Project Structure

```
my-api/
├── cmd/
│   ├── api/
│   │   └── main.go              # Entry point with Fiber v3 config
│   └── service.go               # APIService with DB/Redis clients
├── internal/
│   ├── config/
│   │   └── config.go            # Env config with timeout support
│   ├── handlers/
│   │   ├── handler.go           # Setup + logging middleware
│   │   ├── util/
│   │   │   └── util.go          # Error/success helpers
│   │   └── user/
│   │       └── handler.go       # Domain handler with Handle() method
│   ├── models/                  # SQLC generated models
│   └── repository/
│       └── db.go                # MySQL + Redis connections
├── database/
│   ├── migrations/              # Goose migrations
│   └── queries/                 # SQLC queries
├── docker-compose.yml           # MySQL, Redis, Kafka
├── Dockerfile
├── Makefile
├── sqlc.yaml
└── .env.example
```

## Configuration

Create a `.env` file:

```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=myapi_user
DB_PASSWORD=myapi_pass
DB_NAME=myapi_db

# Server
SERVER_PORT=8080

# Infrastructure
REDIS_HOST=localhost
REDIS_PORT=6379
KAFKA_HOST=localhost
KAFKA_PORT=9092

# Timeouts (Go duration format)
READ_TIMEOUT=5s
WRITE_TIMEOUT=10s
IDLE_TIMEOUT=0s
```

## Writing Handlers

Example domain handler:

```go
package user

import (
	"database/sql"
	
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"myapp/internal/handlers/util"
	"myapp/internal/models"
)

type GetByIDHandler struct {
	DB      *sql.DB
	Queries *models.Queries
	Redis   *redis.Client
}

func (h *GetByIDHandler) Handle(c fiber.Ctx) error {
	// Type-safe URI parameter binding
	var params struct {
		ID int `uri:"id"`
	}
	if err := c.Bind().URI(&params); err != nil {
		return util.HandleError(c, util.BuildErrorWithCode(
			fiber.StatusBadRequest,
			"Invalid ID format",
			"INVALID_ID_FORMAT",
		))
	}
	
	// Pass c directly (implements context.Context)
	user, err := h.Queries.GetUser(c, h.DB, int32(params.ID))
	if err != nil {
		return util.HandleError(c, err)
	}
	
	return c.JSON(util.SuccessResponse{Data: user})
}
```

## Database Workflow

### 1. Create Migration

```bash
make migrate-create
# Enter name: create_users_table
```

### 2. Write SQL Queries

```sql
-- database/queries/user.sql
-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: CreateUser :execresult
INSERT INTO users (name, email) VALUES (?, ?);
```

### 3. Generate Go Code

```bash
make sqlc-generate
```

## Adding Entities

```bash
ready-go add entity Product
```

Creates:
- `internal/entity/product.go` - Struct definition
- `database/migrations/xxx_create_products.sql` - Migration
- `database/queries/product.sql` - SQLC queries

## Error Handling

```go
return util.HandleError(c, util.BuildErrorWithCode(
	fiber.StatusNotFound,
	"User not found",
	"USER_NOT_FOUND",
))
```

Response:
```json
{
	"http_code": 404,
	"message": "User not found",
	"code": "USER_NOT_FOUND"
}
```

## Make Commands

```bash
make docker-up        # Start MySQL, Redis, Kafka
make docker-down      # Stop services
make migrate-up       # Run migrations
make migrate-down     # Rollback migrations
make migrate-create   # Create new migration
make sqlc-generate    # Generate SQLC models
make run-api          # Run dev server
make build-api        # Build binary
```

## Requirements

- Go 1.23+
- Docker & Docker Compose
- sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- goose: `go install github.com/pressly/goose/v3/cmd/goose@latest`

## Tech Stack

- [Fiber v3](https://gofiber.io/) - Web framework
- [MySQL](https://www.mysql.com/) - Database
- [SQLC](https://sqlc.dev/) - Type-safe SQL
- [Goose](https://github.com/pressly/goose) - Migrations
- [Redis](https://redis.io/) - Cache
- [Kafka](https://kafka.apache.org/) - Message queue

## Changelog

We follow semantic versioning and are committed to backwards compatibility.

**Promise:** All future versions will maintain backwards compatibility for generated projects. Your existing projects will continue to work when you regenerate or upgrade the CLI.

### v3.0.1 - Fiber v3 Upgrade
- Upgraded to Fiber v3 with built-in timeout support
- Requires Go 1.23+
- New handler pattern: `GetByIDHandler` with `Handle()` method
- Type-safe parameter binding via `c.Bind().URI()`
- Environment-based timeout configuration (READ_TIMEOUT, WRITE_TIMEOUT, IDLE_TIMEOUT)

### v3.0.0 - Simplified Architecture
- Removed complex layers (DTO, Service, Repository abstractions)
- Direct handler → SQLC model pattern
- Reduced dependencies from 10+ to 4 core deps
- Flag-only CLI interface
- MySQL + Redis + Kafka always included

### v2.x and earlier
- See [releases](https://github.com/muazwzxv/ready-go-cli/releases) for historical changelog

## License

MIT
