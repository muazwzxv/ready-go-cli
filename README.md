# Ready-Go CLI

A CLI tool to scaffold production-ready Go projects with clean architecture, complete with Docker setup, database migrations, and a sample CRUD API.

## Features

- 🏗️ **Clean Architecture**: Entity, DTO, Repository, Service, and Handler layers
- 📁 **Modular Structure**: Handlers and services organized by domain in subdirectories
- 🔌 **Self-Registering Routes**: Handlers register their own routes for better encapsulation
- 🐳 **Docker Ready**: Pre-configured Docker Compose with MySQL, Redis, and Kafka
- 🔄 **Database Migrations**: Built-in migration support with goose
- ⚙️ **Multi-Config**: Environment variables → TOML → Defaults
- 🚀 **Ready to Run**: Generated projects compile and run immediately
- 🎨 **Customizable**: Configurable entity names, module paths, and services
- 📦 **Standalone Binary**: All templates embedded, no external dependencies
- 🧹 **Clean Code**: Generated code is comment-free and production-ready

## What's New (v1.1 - Nov 2025)

### Improved Project Structure
- **Domain-Driven Organization**: Handlers and services are now organized in subdirectories by domain (e.g., `handler/user/`, `service/user/`)
- **One File Per Handler Method**: Each handler method is in its own file for better maintainability
- **Self-Registering Handlers**: Each handler implements `RegisterRoutes(app *fiber.App)` method
- **Structured Dependencies**: `ApplicationContext` now uses `Services` and `Handlers` structs for better organization

### Code Quality
- **Comment-Free Generated Code**: Cleaner output without inline comments
- **Package Aliases**: Better import management with aliases like `healthHandler`, `userHandler`
- **Focused Files**: Single-responsibility files make the codebase easier to navigate

## Installation

### Option 1: Build from source

```bash
cd ready-go-cli
go build -o ready-go cmd/ready-go/main.go

# Install to your $GOBIN or $PATH
mv ready-go $GOPATH/bin/
# or
mv ready-go /usr/local/bin/
```

### Option 2: Install with Go

```bash
go install github.com/muazwzxv/ready-go-cli/cmd/ready-go@latest
```

## Quick Start

```bash
# Create a new project
ready-go new my-api --module github.com/myorg/my-api --sample-api Product

# Navigate to project
cd my-api

# Start services
make up

# Run migrations
make migrate-up

# Start the application
make run
```

Your API will be available at `http://localhost:8080`

## Usage

### Basic Command

```bash
ready-go new <project-name> [flags]
```

### Available Flags

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `--module` | `-m` | `github.com/username/<project-name>` | Go module name |
| `--description` | `-d` | Auto-generated | Project description |
| `--author` | | Empty | Author name |
| `--output` | `-o` | `.` | Output directory |
| `--sample-api` | | `User` | Sample API entity name (e.g., User, Product, Order) |
| `--with-redis` | | `true` | Include Redis in docker-compose |
| `--with-kafka` | | `true` | Include Kafka in docker-compose |
| `--skip-git` | | `false` | Skip git initialization |
| `--interactive` | `-i` | `false` | Interactive mode with prompts |

### Examples

**Create a simple user management API:**
```bash
ready-go new user-service --module github.com/mycompany/user-service
```

**Create a product catalog API without Kafka:**
```bash
ready-go new product-catalog \
  --module github.com/shop/products \
  --sample-api Product \
  --with-kafka=false
```

**Interactive mode:**
```bash
ready-go new my-project -i
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
│   │   └── config.go           # Multi-config loader
│   ├── database/
│   │   ├── database.go         # MySQL connection
│   │   └── migrations/         # SQL migrations
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
│   │   └── {entity}_repository.go
│   ├── service/                # Business logic layer
│   │   └── {entity}/           # Organized by domain
│   │       ├── {entity}.go     # Service interface
│   │       └── create_{entity}_service.go
│   └── handler/                # HTTP handlers
│       ├── health/             # Health check handlers
│       │   └── health_handler.go
│       ├── {entity}/           # Entity-specific handlers
│       │   ├── {entity}_handler.go
│       │   └── create_{entity}_handler.go
│       └── middleware.go       # Common middleware
├── Dockerfile                  # Multi-stage Docker build
├── docker-compose.yml          # MySQL, Redis, Kafka setup
├── Makefile                    # Dev commands
├── config.toml                 # Configuration file
├── .env.docker                 # Docker environment
├── .env.example                # Example environment
├── .gitignore
└── README.md
```

### Project Structure Highlights

- **Modular Organization**: Handlers and services are organized by domain (e.g., `handler/user/`, `service/user/`)
- **Self-Registering Handlers**: Each handler package implements `RegisterRoutes(app *fiber.App)` for clean route management
- **Structured Dependencies**: `ApplicationContext` organizes dependencies into `Services` and `Handlers` structs
- **Single Responsibility**: Each handler method is in its own file for better maintainability
- **Clean Imports**: Package aliases prevent naming conflicts (e.g., `healthHandler`, `userHandler`)

## Generated API Endpoints

The CLI generates a foundation CRUD API for your sample entity. Currently includes:

**Entity Endpoints:**
- `POST /api/v1/{entity}s` - Create new entity

**Health Check Endpoints:**
- `GET /health` - Application health with database status
- `GET /health/ready` - Readiness check for orchestrators

### Easy to Extend

The generated structure makes it simple to add more endpoints:

```go
// In handler/{entity}/{entity}_handler.go
func (h *UserHandler) RegisterRoutes(app *fiber.App) {
    app.Post("/api/v1/users", h.CreateUser)
    // Add more routes here:
    app.Get("/api/v1/users/:id", h.GetUser)
    app.Put("/api/v1/users/:id", h.UpdateUser)
    app.Delete("/api/v1/users/:id", h.DeleteUser)
}
```

Each handler method can be in its own file (e.g., `get_user_handler.go`, `update_user_handler.go`) for better organization.

## Working with Generated Projects

### Available Make Commands

```bash
make up              # Start all Docker services
make down            # Stop all Docker services
make migrate-up      # Run database migrations
make migrate-down    # Rollback database migrations
make run             # Run the application
make build           # Build the binary
make test            # Run tests
make lint            # Run linter
make clean           # Clean build artifacts
```

### Configuration

Generated projects support three configuration sources (in order of precedence):

1. **Environment Variables** (highest priority)
2. **config.toml file**
3. **Default values** (lowest priority)

Example configuration:
```bash
# Using environment variables
export DB_HOST=localhost
export DB_PORT=3306
export APP_PORT=8080

# Or edit config.toml
vim config.toml

# Or use .env files with Docker
cp .env.example .env.docker
```

### Database Migrations

```bash
# Create a new migration
goose -dir internal/database/migrations create my_migration sql

# Run migrations
make migrate-up

# Rollback last migration
make migrate-down
```

## Architecture

The generated projects follow clean architecture principles with a modular, domain-driven structure:

### Layer Organization

1. **Entity Layer**: Core business models and domain logic
   - Pure Go structs with business methods
   - No external dependencies

2. **DTO Layer**: Data transfer objects for API requests/responses
   - Separate request and response types
   - Validation-ready structures

3. **Repository Layer**: Database operations and data access
   - Interface-based design
   - MySQL implementation with sqlx

4. **Service Layer**: Business logic and orchestration
   - Organized by domain (e.g., `service/user/`)
   - Interface definitions in `{entity}.go`
   - Implementation in separate files by operation

5. **Handler Layer**: HTTP request handling and routing
   - Organized by domain (e.g., `handler/user/`)
   - Self-registering routes via `RegisterRoutes()`
   - One handler method per file for clarity

### Dependency Flow

```
Handler → Service → Repository → Entity
   ↓         ↓          ↓
  DTO    Interface   Database
```

### Application Wiring

The `ApplicationContext` in `internal/application.go` manages dependencies:

```go
type ApplicationContext struct {
    Services *Services  // All business logic services
    Handlers *Handlers  // All HTTP handlers
}

type Services struct {
    UserService service.UserService
}

type Handlers struct {
    UserHandler   *userHandler.UserHandler
    HealthHandler *healthHandler.HealthHandler
}
```

This structure makes it easy to:
- Add new services and handlers
- Mock dependencies for testing
- Maintain clear separation of concerns

## Migration Guide (v1.0 → v1.1)

If you have an existing project generated with v1.0, here's what changed:

### Structure Changes

**Old Structure (v1.0):**
```
internal/
  handler/
    health_handler.go
    user_handler.go
  service/
    interfaces.go
    user_service.go
```

**New Structure (v1.1):**
```
internal/
  handler/
    health/
      health_handler.go
    user/
      user_handler.go
      create_user_handler.go
  service/
    user/
      user.go
      create_user_service.go
```

### Key Changes

1. **Handlers in Subdirectories**: Move each handler to its own subdirectory
2. **Self-Registering Routes**: Add `RegisterRoutes(app *fiber.App)` method to handlers
3. **Structured ApplicationContext**: Use `Services` and `Handlers` structs
4. **Package Aliases**: Import handlers with aliases (e.g., `userHandler "path/to/handler/user"`)

### Migration Steps

For new projects, simply regenerate with v1.1. For existing projects, you can:
1. Keep your v1.0 structure (still works fine)
2. Manually refactor to match v1.1 structure (recommended for long-term maintainability)
3. Create a new v1.1 project and port your business logic

## Requirements

- Go 1.23 or later
- Docker and Docker Compose (for generated projects)
- Make (optional, but recommended)

## Development

To modify the CLI itself:

```bash
# Clone the repository
git clone https://github.com/muazwzxv/ready-go-cli.git
cd ready-go-cli

# Make changes to templates in cmd/ready-go/templates/

# Rebuild
go build -o ready-go cmd/ready-go/main.go

# Test
./ready-go new test-project
```

## Template Customization

Templates are embedded in the binary but can be customized:

1. Templates are located in `cmd/ready-go/templates/`
2. Modify any `.tmpl` file
3. Rebuild the CLI: `go build -o ready-go cmd/ready-go/main.go`

Available template variables:
- `{{.ProjectName}}` - Project name
- `{{.ModuleName}}` - Go module name
- `{{.SampleAPIName}}` - Entity name (e.g., "User", "Product")
- `{{.SampleAPINameLower}}` - Lowercase entity name (e.g., "user", "product")
- `{{.SampleTableName}}` - Table name (e.g., "users", "products")
- `{{.AppPort}}`, `{{.MySQLPort}}`, etc. - Port configurations

## Troubleshooting

**Q: CLI can't find templates**
A: Make sure the binary was built from the correct location. Templates are embedded during build.

**Q: Generated project doesn't compile**
A: Run `go mod tidy` in the generated project directory to ensure all dependencies are downloaded.

**Q: Docker services won't start**
A: Check if ports are already in use: `lsof -i :3306` (MySQL), `lsof -i :6379` (Redis), etc.

**Q: Migrations fail**
A: Ensure MySQL is running and accessible. Check connection settings in `config.toml` or environment variables.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Author

Ready-Go CLI Team

---

**Happy coding! 🚀**
