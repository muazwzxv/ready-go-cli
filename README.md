# Ready-Go CLI

A CLI tool to scaffold production-ready Go projects with clean architecture, complete with Docker setup, database migrations, and type-safe SQL queries.

**Latest Version**: v1.1.0

## Features

- **Clean Architecture**: Entity, DTO, Repository, Service, and Handler layers
- **Modular Structure**: Handlers and services organized by domain in subdirectories
- **Self-Registering Routes**: Handlers register their own routes for better encapsulation
- **Docker Ready**: Pre-configured Docker Compose with MySQL, Redis, and Kafka
- **Database Migrations**: Built-in migration support with goose
- **Type-Safe SQL**: Integrated sqlc for compile-time SQL query validation
- **Viper Configuration**: Industry-standard config management with automatic environment variable override
- **Multi-Source Config**: Defaults → TOML file → environment variables (priority order)
- **Ready to Run**: Generated projects compile and run immediately
- **Customizable**: Configurable entity names, module paths, and services
- **Standalone Binary**: All templates embedded, no external dependencies

## Installation

### Option 1: Install with Go

```bash
go install github.com/muazwzxv/ready-go-cli/cmd/ready-go@latest
```

### Option 2: Build from source

```bash
git clone https://github.com/muazwzxv/ready-go-cli.git
cd ready-go-cli
go build -o ready-go ./cmd/ready-go

# Move to your PATH
mv ready-go ~/go/bin/
# or
mv ready-go /usr/local/bin/
```

## Quick Start

```bash
# Create a new project
ready-go new my-api

# Navigate to project
cd my-api

# Start services with Docker
make up

# Run migrations
make migrate-up

# Generate type-safe SQL code
make sqlc-generate

# Start the application
make run
```

Your API will be available at `http://localhost:8080`

## Usage

### Command Syntax

```bash
ready-go new [options] <project-name>
```

The CLI supports two modes:

1. **Interactive mode** (default): Prompts you for all configuration options
2. **Flag mode**: Specify options directly via command-line flags

### Available Flags

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `--module` | `-m` | `github.com/username/<project-name>` | Go module path |
| `--description` | `-d` | Auto-generated | Project description |
| `--author` | | Empty | Author name |
| `--output` | `-o` | `.` | Output directory |
| `--sample-api` | | `User` | Sample API entity name (e.g., User, Product, Order) |
| `--with-redis` | | `true` | Include Redis in docker-compose |
| `--with-kafka` | | `true` | Include Kafka in docker-compose |
| `--skip-git` | | `false` | Skip git initialization |
| `--interactive` | `-i` | `false` | Force interactive mode with prompts |

### Interactive Prompts

When you run `ready-go new` in interactive mode (default), you'll be prompted for:

| Prompt | Description | Default |
|--------|-------------|---------|
| **Project Name** | Name of your project | Required argument |
| **Module Name** | Go module path | `github.com/username/<project-name>` |
| **Description** | Project description | Auto-generated |
| **Author** | Author name | Empty |
| **Sample API Entity** | Entity name for sample CRUD API | `User` |
| **Table Name** | Database table name | Auto-pluralized from entity |
| **Application Port** | HTTP server port | `8080` |
| **MySQL Port** | MySQL database port | `3306` |
| **Include Redis** | Add Redis to docker-compose | `yes` |
| **Include Kafka** | Add Kafka to docker-compose | `yes` |

### Examples

#### Interactive Mode (Default)

**Basic project with interactive prompts:**
```bash
ready-go new my-api
# CLI will prompt you for:
# - Module name
# - Description
# - Author
# - Sample API entity
# - Table name
# - Ports
# - Include Redis/Kafka
```

**Force interactive mode explicitly:**
```bash
ready-go new --interactive my-api
# or
ready-go new -i my-api
```

#### Using Command-Line Flags

**Simple project with custom module:**
```bash
ready-go new --module github.com/mycompany/user-api user-api
```

**Product catalog without Kafka:**
```bash
ready-go new \
  --module github.com/mycompany/product-catalog \
  --sample-api Product \
  --with-kafka=false \
  product-catalog
```

**Full configuration with flags:**
```bash
ready-go new \
  --module github.com/mycompany/order-service \
  --description "Microservice for order management" \
  --author "John Doe" \
  --sample-api Order \
  --with-redis=true \
  --with-kafka=false \
  order-service
```

**Minimal setup (no Redis, no Kafka):**
```bash
ready-go new \
  --module github.com/mycompany/minimal-api \
  --with-redis=false \
  --with-kafka=false \
  minimal-api
```

#### Mixed Mode (Flags + Interactive)

**Specify some options via flags, get prompted for the rest:**
```bash
ready-go new --module github.com/mycompany/my-api my-api
# CLI will prompt for remaining options:
# - Description
# - Author
# - Sample API entity
# - etc.
```

## Generated Project Structure

```
my-project/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── application.go           # App initialization & wiring
│   ├── config/
│   │   └── config.go           # Multi-source config loader
│   ├── database/
│   │   ├── database.go         # MySQL connection setup
│   │   ├── sqlc.yaml           # sqlc configuration
│   │   ├── query/              # SQL queries for sqlc
│   │   │   └── sample.sql
│   │   ├── store/              # Generated sqlc code (auto-generated)
│   │   └── migrations/         # SQL migration files
│   │       └── 001_create_samples_table.sql
│   ├── entity/                 # Domain models
│   │   └── {entity}.go
│   ├── dto/
│   │   ├── request/            # API request DTOs
│   │   │   ├── common.go
│   │   │   └── {entity}_request.go
│   │   └── response/           # API response DTOs
│   │       ├── common.go
│   │       ├── error_response.go
│   │       └── {entity}_response.go
│   ├── repository/             # Data access layer
│   │   ├── interfaces.go
│   │   ├── errors.go
│   │   └── {entity}_repository.go
│   ├── service/                # Business logic layer
│   │   └── {entity}/           # Organized by domain
│   │       ├── {entity}.go
│   │       └── create_{entity}_service.go
│   └── handler/                # HTTP handlers
│       ├── health/
│       │   └── health_handler.go
│       ├── {entity}/
│       │   ├── {entity}_handler.go
│       │   └── create_{entity}_handler.go
│       └── middleware.go
├── Dockerfile                  # Multi-stage Docker build
├── docker-compose.yml          # MySQL, Redis, Kafka setup
├── Makefile                    # Development commands
├── config.toml                 # Configuration file
├── .env.docker                 # Docker environment
├── .env.example                # Example environment variables
├── .gitignore
└── README.md
```

## Tools & Technologies Included

### Out of the Box

- **Web Framework**: [Fiber v2](https://gofiber.io/) - Fast HTTP framework
- **Database**: MySQL 8.0 with [sqlx](https://github.com/jmoiron/sqlx)
- **Type-Safe SQL**: [sqlc](https://sqlc.dev/) - Generate Go code from SQL
- **Configuration**: [Viper](https://github.com/spf13/viper) - Flexible configuration with automatic env var binding
- **Migrations**: [goose](https://github.com/pressly/goose) - Database migration tool
- **Containerization**: Docker & Docker Compose
- **Optional Services**: Redis 7.2, Kafka 3.5 with UI

### Architecture Pattern

Generated projects follow **Clean Architecture** principles:

```
Handler → Service → Repository → Database
   ↓         ↓          ↓
  DTO    Interface   Entity
```

**Benefits:**
- Testable business logic (services are pure Go)
- Swappable implementations (interface-based repositories)
- Clear separation of concerns
- Easy to mock for unit tests

## Working with Generated Projects

### Make Commands

```bash
make help             # Show all available commands
make up               # Start all Docker services
make down             # Stop all Docker services
make logs             # View service logs
make migrate-up       # Run database migrations
make migrate-down     # Rollback last migration
make migrate-status   # Check migration status
make migrate-create   # Create new migration (use NAME=migration_name)
make sqlc-generate    # Generate Go code from SQL queries
make run              # Run the application locally
make build            # Build the binary
make test             # Run tests
make clean            # Clean build artifacts and stop services
```

### Configuration

Generated projects use **Viper** for flexible, multi-source configuration with the following priority:

1. **Environment Variables** (highest priority)
2. **config.toml file**
3. **Default values** (lowest priority)

#### Environment Variables

Environment variables automatically override config file values:

```bash
# Server configuration
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080
export SERVER_READ_TIMEOUT=5s
export SERVER_WRITE_TIMEOUT=10s

# Database configuration
export DATABASE_HOST=localhost
export DATABASE_PORT=3306
export DATABASE_USER=myuser
export DATABASE_PASSWORD=mypassword
export DATABASE_DATABASE=myapp
export DATABASE_MAX_OPEN_CONNS=50
export DATABASE_MAX_IDLE_CONNS=20
```

**Naming convention**: Nested keys are flattened with underscores
- `server.host` → `SERVER_HOST`
- `database.max_open_conns` → `DATABASE_MAX_OPEN_CONNS`

#### TOML Configuration File

Edit `config.toml` for structured configuration:

```toml
[server]
host = "0.0.0.0"
port = 8080
read_timeout = "5s"
write_timeout = "10s"

[database]
host = "localhost"
port = 3306
user = "appuser"
password = "apppassword"
database = "myapp"
max_open_conns = 25
max_idle_conns = 10
conn_max_lifetime = "5m"
```

### Database Operations

#### Migrations

```bash
# Create a new migration
make migrate-create NAME=add_users_table

# Run all pending migrations
make migrate-up

# Rollback the last migration
make migrate-down

# Check migration status
make migrate-status
```

Migration files are created in `internal/database/migrations/` with timestamp prefixes.

#### sqlc - Type-Safe SQL Queries

The generated project uses [sqlc](https://sqlc.dev/) to generate type-safe Go code from SQL queries.

**Workflow:**

1. Write SQL queries in `internal/database/query/*.sql`:

```sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: CreateUser :execresult
INSERT INTO users (name, email, status, created_at, updated_at)
VALUES (?, ?, ?, NOW(), NOW());

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
```

2. Generate Go code:

```bash
make sqlc-generate
```

3. Use the generated type-safe functions in your repository:

```go
import "your-module/internal/database/store"

queries := store.New()
user, err := queries.GetUserByID(ctx, db, userID)
users, err := queries.ListUsers(ctx, db, store.ListUsersParams{
    Limit:  20,
    Offset: 0,
})
```

**Benefits:**
- Compile-time SQL validation
- Type-safe database operations
- Auto-generated Go structs from schema
- No runtime reflection overhead
- Catches SQL errors at build time

### Generated API Endpoints

**Health Checks:**
- `GET /health` - Basic health check with database status
- `GET /health/ready` - Readiness check for Kubernetes/orchestrators

**Sample Entity CRUD:**
- `POST /api/v1/{entities}` - Create new entity

The generated code provides a foundation. You can easily extend it by adding more handler methods:

```go
// In handler/{entity}/{entity}_handler.go
func (h *UserHandler) RegisterRoutes(app *fiber.App) {
    app.Post("/api/v1/users", h.CreateUser)
    // Add more routes:
    app.Get("/api/v1/users", h.ListUsers)
    app.Get("/api/v1/users/:id", h.GetUser)
    app.Put("/api/v1/users/:id", h.UpdateUser)
    app.Delete("/api/v1/users/:id", h.DeleteUser)
}
```

## Requirements

### CLI Tool
- Go 1.23 or later

### Generated Projects
- Go 1.23 or later
- Docker and Docker Compose (for containerized development)
- Make (optional but recommended)
- sqlc (for SQL code generation) - `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

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
make up && make migrate-up && make run
```

### Template Variables

Templates are located in `cmd/ready-go/templates/` and support these variables:

- `{{.ProjectName}}` - Project name
- `{{.ModuleName}}` - Go module name
- `{{.Description}}` - Project description
- `{{.Author}}` - Author name
- `{{.SampleAPIName}}` - Entity name (e.g., "User", "Product")
- `{{.SampleAPINameLower}}` - Lowercase entity name (e.g., "user", "product")
- `{{.SampleTableName}}` - Table name (e.g., "users", "products")
- `{{.AppPort}}` - Application HTTP port
- `{{.MySQLPort}}` - MySQL port
- `{{.GoVersion}}` - Go version requirement
- `{{.WithRedis}}` - Include Redis (boolean)
- `{{.WithKafka}}` - Include Kafka (boolean)

## Troubleshooting

**Q: Generated project doesn't compile**

A: Run `go mod tidy` in the generated project directory. If using sqlc, also run `make sqlc-generate` to generate the database code.

**Q: Docker services won't start**

A: Check if ports are already in use:
```bash
lsof -i :8080   # Application
lsof -i :3306   # MySQL
lsof -i :6379   # Redis
lsof -i :9092   # Kafka
```

**Q: Migrations fail**

A: Ensure MySQL is running (`docker-compose ps`) and check connection settings in `config.toml` or environment variables.

**Q: sqlc generation fails**

A: Make sure sqlc is installed: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

**Q: "no Go files" error when building**

A: The main.go file is in `cmd/server/`. Build with: `go build -o app ./cmd/server`

## What's New in v1.1.0

### Viper Configuration Integration

This release introduces **Viper** for configuration management, replacing the previous manual implementation:

**Key Improvements:**
- ✨ **Simpler Code**: Configuration code reduced by ~60% (106 lines vs 266 lines)
- 🔄 **Automatic Env Var Binding**: No more manual environment variable parsing
- 🐳 **Better Docker Support**: Environment variables now work seamlessly in containers
- 🎯 **Default Values**: All config keys have sensible defaults
- 📦 **Industry Standard**: Using the widely-adopted Viper library

**Breaking Change - Environment Variable Naming:**

Environment variables now follow a consistent `DATABASE_*` naming convention:

```bash
# ❌ Old (v1.0.0)
DB_HOST=localhost
DB_PORT=3306
DB_USER=appuser

# ✅ New (v1.1.0)
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_USER=appuser
```

**Migration:** Existing projects can update their `.env` files to use the new naming, or regenerate with the latest CLI version. See [CHANGELOG.md](CHANGELOG.md) for detailed migration instructions.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License - see LICENSE file for details

## Author

Muaz - [@muazwzxv](https://github.com/muazwzxv)

---

**Built with Go. Ready to scale.**
