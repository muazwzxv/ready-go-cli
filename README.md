# Ready-Go CLI

A CLI tool to scaffold production-ready Go projects with a simplified, practical architecture.

**Latest Version**: v3.0.0

## Features

- **Simplified Architecture**: Handlers directly use SQLC models - no complex abstraction layers
- **Docker Ready**: Pre-configured Docker Compose with MySQL, Redis, and Kafka
- **Database Migrations**: Built-in migration support with goose
- **Type-Safe SQL**: Integrated sqlc for compile-time SQL query validation
- **Environment Configuration**: Simple dotenv-based configuration
- **Request Logging**: Built-in HTTP request logging with slog
- **Error Handling**: Structured error responses with error codes
- **Domain Organization**: Handlers organized by domain in subdirectories
- **Standalone Binary**: All templates embedded, no external dependencies

## Installation

### Build from source

```bash
git clone https://github.com/muazwzxv/ready-go-cli.git
cd ready-go-cli
go build -o ready-go ./cmd/ready-go

# Move to your PATH
mv ready-go ~/go/bin/
```

## Quick Start

```bash
# Create a new project
ready-go new my-api --module github.com/mycompany/my-api

# Navigate to project
cd my-api

# Start services with Docker
make docker-up

# Run migrations
make migrate-up

# Generate type-safe SQL code
make sqlc-generate

# Start the application
make run-api
```

Your API will be available at `http://localhost:8080`

## Usage

### Command Syntax

```bash
ready-go new [flags] <project-name>
```

### Available Flags

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `--module` | `-m` | `github.com/username/<project-name>` | Go module path |
| `--port` | | `8080` | Server HTTP port |
| `--db-port` | | `3306` | MySQL port |
| `--redis-port` | | `6379` | Redis port |
| `--kafka-port` | | `9092` | Kafka port |
| `--sample-name` | | `User` | Sample entity name (e.g., User, Product) |

### Examples

**Simple project:**
```bash
ready-go new my-api
```

**With custom module and entity:**
```bash
ready-go new product-service \
  --module github.com/mycompany/product-service \
  --sample-name Product
```

**Custom ports:**
```bash
ready-go new my-api \
  --port 3000 \
  --db-port 3307 \
  --redis-port 6380
```

## Generated Project Structure

```
my-project/
├── cmd/
│   ├── api/
│   │   └── main.go              # Application entry point
│   └── service.go               # APIService struct
├── internal/
│   ├── config/
│   │   └── config.go           # Environment config loader
│   ├── handlers/
│   │   ├── handler.go          # SetupHandler + middleware
│   │   └── user/               # Domain-specific handlers
│   │       └── handler.go
│   ├── response/
│   │   └── response.go         # Error/success responses
│   ├── models/                 # SQLC generated code
│   └── repository/
│       └── db.go               # MySQL + Redis connections
├── database/
│   ├── migrations/             # Goose migration files
│   └── queries/                # SQLC query files
├── Dockerfile                  # Multi-stage Docker build
├── docker-compose.yml          # MySQL, Redis, Kafka
├── Makefile                    # Development commands
├── sqlc.yaml                   # SQLC configuration
├── .env.example                # Example environment variables
└── README.md
```

## Tools & Technologies

### Out of the Box

- **Web Framework**: [Fiber v2](https://gofiber.io/) - Fast HTTP framework
- **Database**: MySQL 8.0 with standard `database/sql`
- **Type-Safe SQL**: [sqlc](https://sqlc.dev/) - Generate Go code from SQL
- **Configuration**: [godotenv](https://github.com/joho/godotenv) - Simple env file loading
- **Migrations**: [goose](https://github.com/pressly/goose) - Database migration tool
- **Cache**: Redis 7.x with [go-redis](https://github.com/redis/go-redis)
- **Message Queue**: Kafka (optional, ready to use)
- **Logging**: Go standard library `log/slog`

### Architecture

Generated projects use a practical, simplified architecture:

```
HTTP Request → Handler → SQLC Models → Database
                   ↓
            Response (JSON)
```

**Benefits:**
- Minimal boilerplate code
- Direct database access via SQLC
- Type-safe SQL queries at compile time
- Easy to understand and extend
- Ready for monolithic growth with domain subfolders

## Working with Generated Projects

### Make Commands

```bash
make docker-up        # Start Docker services (MySQL, Redis, Kafka)
make docker-down      # Stop Docker services
make migrate-up       # Run database migrations
make migrate-down     # Rollback last migration
make migrate-create   # Create new migration
make sqlc-generate    # Generate Go code from SQL queries
make run-api          # Run the application locally
make build-api        # Build the binary
```

### Configuration

Configuration is loaded from environment variables via `.env` file:

```bash
# Copy example file
cp .env.example .env

# Edit as needed
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
```

### Database Operations

#### Migrations

```bash
# Create a new migration
make migrate-create
# Enter migration name when prompted

# Run all pending migrations
make migrate-up

# Rollback the last migration
make migrate-down
```

Migration files are created in `database/migrations/` with timestamp prefixes.

#### sqlc - Type-Safe SQL Queries

**Workflow:**

1. Write SQL queries in `database/queries/*.sql`:

```sql
-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: CreateUser :execresult
INSERT INTO users (name, email, created_at)
VALUES (?, ?, NOW());

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;
```

2. Generate Go code:

```bash
make sqlc-generate
```

3. Use in handlers:

```go
func (h *Handler) GetByID(c *fiber.Ctx) error {
    id := c.ParamsInt("id")
    user, err := h.Queries.GetUser(c.Context(), h.DB, int32(id))
    if err != nil {
        return response.HandleError(c, err)
    }
    return c.JSON(response.SuccessResponse{Data: user})
}
```

### Adding New Entities

Use the `add entity` command to scaffold new entities:

```bash
cd my-project
ready-go add entity Product
```

This creates:
- `internal/entity/product.go` - Entity struct
- `database/migrations/xxx_create_products.sql` - Migration
- `database/queries/product.sql` - SQLC queries

### Error Handling

Return structured errors with codes:

```go
return response.HandleError(c, response.BuildErrorWithCode(
    fiber.StatusNotFound,
    "User not found",
    "USER_NOT_FOUND",
))
```

Response format:
```json
{
  "http_code": 404,
  "message": "User not found",
  "code": "USER_NOT_FOUND"
}
```

## Requirements

### CLI Tool
- Go 1.21 or later

### Generated Projects
- Go 1.21 or later
- Docker and Docker Compose
- Make (optional but recommended)
- sqlc - `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- goose - `go install github.com/pressly/goose/v3/cmd/goose@latest`

## Development

### Modifying the CLI

```bash
# Clone the repository
git clone https://github.com/muazwzxv/ready-go-cli.git
cd ready-go-cli

# Make changes to templates in cmd/ready-go/templates/

# Rebuild the CLI
go build -o ready-go ./cmd/ready-go

# Test with a new project
./ready-go new test-project
cd test-project
make docker-up && make migrate-up && make run-api
```

## What's New in v3.0.0

### Simplified Architecture

This major release removes complexity while keeping productivity:

**Removed:**
- ❌ Dependency injection (samber/do)
- ❌ Viper configuration system
- ❌ Multi-layer architecture (DTO → Service → Repository)
- ❌ Interactive prompts
- ❌ Optional Redis/Kafka flags (now always included)

**Changed:**
- ✅ Handlers directly use SQLC models
- ✅ Simple godotenv configuration
- ✅ Flattened project structure
- ✅ Flag-only CLI interface
- ✅ Always includes MySQL + Redis + Kafka

**Benefits:**
- 53% fewer template files (34 → 16)
- Reduced dependencies (10+ → 4)
- Less boilerplate code
- Easier to understand and extend
- Faster project generation

### Migration from v2.x

Projects generated with v2.x will continue to work. New projects use the simplified v3.0 structure.

To upgrade existing projects:
1. Keep your existing project as-is
2. Generate a new project with v3.0
3. Copy your business logic to the new structure

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing-feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see LICENSE file for details

## Author

Muaz - [@muazwzxv](https://github.com/muazwzxv)

---

**Built with Go. Ready to code.**
