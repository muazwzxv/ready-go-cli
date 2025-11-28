# Ready-Go CLI

A CLI tool to scaffold production-ready Go projects with clean architecture, complete with Docker setup, database migrations, and a sample CRUD API.

## Features

- 🏗️ **Clean Architecture**: Entity, DTO, Repository, Service, and Handler layers
- 🐳 **Docker Ready**: Pre-configured Docker Compose with MySQL, Redis, and Kafka
- 🔄 **Database Migrations**: Built-in migration support with goose
- ⚙️ **Multi-Config**: Environment variables → TOML → Defaults
- 🚀 **Ready to Run**: Generated projects compile and run immediately
- 🎨 **Customizable**: Configurable entity names, module paths, and services
- 📦 **Standalone Binary**: All templates embedded, no external dependencies

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
│   ├── dto/
│   │   ├── request/            # API request DTOs
│   │   └── response/           # API response DTOs
│   ├── repository/             # Data access layer
│   ├── service/                # Business logic layer
│   └── handler/                # HTTP handlers
├── Dockerfile                  # Multi-stage Docker build
├── docker-compose.yml          # MySQL, Redis, Kafka setup
├── Makefile                    # Dev commands
├── config.toml                 # Configuration file
├── .env.docker                 # Docker environment
├── .env.example                # Example environment
├── .gitignore
└── README.md
```

## Generated API Endpoints

The CLI generates a complete CRUD API for your sample entity:

- `POST /api/v1/{entity}s` - Create
- `GET /api/v1/{entity}s/:id` - Get by ID
- `PUT /api/v1/{entity}s/:id` - Update
- `DELETE /api/v1/{entity}s/:id` - Delete
- `GET /api/v1/{entity}s` - List with pagination/filtering/sorting
- `POST /api/v1/{entity}s/bulk-update-status` - Bulk operations
- `GET /api/v1/{entity}s/stats` - Statistics

**Health Check Endpoints:**
- `GET /health` - Application health
- `GET /ready` - Readiness check
- `GET /live` - Liveness check

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

The generated projects follow clean architecture principles:

1. **Entity Layer**: Core business models and domain logic
2. **DTO Layer**: Data transfer objects for API requests/responses
3. **Repository Layer**: Database operations and data access
4. **Service Layer**: Business logic and orchestration
5. **Handler Layer**: HTTP request handling and routing

Dependencies flow inward:
```
Handler → Service → Repository → Entity
```

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
