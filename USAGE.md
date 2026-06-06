# Ready-Go CLI - Usage Guide

**Version**: v3.x.x

A command-line tool to scaffold production-ready Go projects with Fiber v3 or Chi, SQLC, and clean architecture.

## Installation

```bash
# Install from source
go install github.com/muazwzxv/ready-go-cli/cmd/ready-go@latest

# Or clone and build
git clone https://github.com/muazwzxv/ready-go-cli.git
cd ready-go-cli
go build -o ready-go ./cmd/ready-go
```

## Quick Start

### Fiber (default)

```bash
ready-go new my-api --module github.com/mycompany/my-api
cd my-api
make docker-up
make migrate-up
make sqlc-generate
make run-api
```

### Chi

```bash
ready-go new my-api --module github.com/mycompany/my-api --router=chi
cd my-api
make docker-up
make migrate-up
make sqlc-generate
make run-api
```

## Commands

### `ready-go new <project-name> [flags]`

Scaffold a new Go project.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--module` | `-m` | `github.com/username/<name>` | Go module path |
| `--router` | `-r` | `fiber` | HTTP router: `fiber` or `chi` |
| `--port` | | `8080` | Server port |
| `--db-port` | | `3306` | MySQL port |
| `--redis-port` | | `6379` | Redis port |
| `--kafka-port` | | `9092` | Kafka port |
| `--sample-name` | | `User` | Sample entity name |

### `ready-go add entity <EntityName>`

Add a new entity to an existing project. Works identically for both Fiber and Chi projects.

```bash
ready-go add entity Product
```

Creates:
- `internal/entity/product.go`
- `database/migrations/xxx_create_products.sql`
- `database/queries/product.sql`

## Generated Project Structure

```
my-api/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── service.go
├── internal/
│   ├── config/
│   ├── handlers/
│   ├── models/
│   └── repository/
├── database/
│   ├── migrations/
│   └── queries/
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── sqlc.yaml
└── .env.example
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

## Configuration

Projects use environment variables (with `.env` file support via godotenv):

```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=myapi_user
DB_PASSWORD=myapi_pass
DB_NAME=myapi_db
SERVER_PORT=8080
REDIS_HOST=localhost
REDIS_PORT=6379
KAFKA_HOST=localhost
KAFKA_PORT=9092
READ_TIMEOUT=5s
WRITE_TIMEOUT=10s
IDLE_TIMEOUT=0s
```

## Router Choice

| | Fiber (default) | Chi |
|---|---|---|
| Handler signature | `func(c fiber.Ctx) error` | `func(w http.ResponseWriter, r *http.Request)` |
| JSON helper | `c.JSON(v)` | `util.WriteJSON(w, code, v)` |
| Param binding | `c.Bind().URI(&params)` | `chi.URLParam(r, "id")` + parse |
| Middleware | `fiber.Handler` | `func(next http.Handler) http.Handler` |
| Dependencies | More | Minimal (chi + stdlib) |

## Requirements

- Go 1.23+
- Docker & Docker Compose
- sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- goose: `go install github.com/pressly/goose/v3/cmd/goose@latest`

## Getting Help

```bash
ready-go --help
ready-go new --help
ready-go add entity --help
```

---

**Happy coding!** 🚀
