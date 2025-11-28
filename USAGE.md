# Ready-Go CLI - Quick Start Guide

A command-line tool to instantly create production-ready Go projects with clean architecture.

## Installation

The `ready-go` binary is already installed at:
```bash
~/go/bin/ready-go
```

Make sure `~/go/bin` is in your PATH.

## Quick Start

Create a new project in 3 steps:

```bash
# 1. Create project
ready-go new my-api --module github.com/mycompany/my-api --sample-api Product

# 2. Navigate and start services
cd my-api && make up

# 3. Run migrations and start app
make migrate-up && make run
```

Your API is now running at `http://localhost:8080`! 🚀

## Basic Usage

### Create a Simple Project
```bash
ready-go new blog-api
```
This creates a project with default settings:
- Module: `github.com/username/blog-api`
- Sample API: User CRUD operations
- Includes: MySQL, Redis, Kafka

### Create with Custom Settings
```bash
ready-go new shop-api \
  --module github.com/mycompany/shop \
  --sample-api Product \
  --author "Your Name"
```

### Create Without Kafka/Redis
```bash
ready-go new simple-api \
  --module github.com/me/simple-api \
  --with-kafka=false \
  --with-redis=false
```

## Command Options

```bash
ready-go new <project-name> [options]
```

### Available Options

| Option | Short | Default | Description |
|--------|-------|---------|-------------|
| `--module` | `-m` | `github.com/username/<name>` | Your Go module path |
| `--sample-api` | | `User` | Name of entity (Product, Order, etc.) |
| `--author` | | Empty | Your name for documentation |
| `--description` | `-d` | Auto-generated | Project description |
| `--output` | `-o` | `.` | Where to create the project |
| `--with-redis` | | `true` | Include Redis in Docker |
| `--with-kafka` | | `true` | Include Kafka in Docker |
| `--skip-git` | | `false` | Don't initialize git repo |
| `--interactive` | `-i` | `false` | Interactive setup mode |

## What You Get

Every generated project includes:

### 🏗️ Clean Architecture
- **Entity Layer**: Domain models with business logic
- **DTO Layer**: Request/response objects
- **Repository Layer**: Database operations
- **Service Layer**: Business logic
- **Handler Layer**: HTTP API endpoints

### 🐳 Docker Environment
- MySQL 8.0 database
- Redis cache (optional)
- Kafka + UI (optional)
- Pre-configured docker-compose.yml

### 🔧 Development Tools
- **Makefile** with common commands
- **Database migrations** with goose
- **Configuration** via ENV/TOML/Defaults
- **Health checks** at /health, /ready, /live

### 📝 API Endpoints

For a "Product" entity, you automatically get:

```
POST   /api/v1/products              # Create product
GET    /api/v1/products/:id          # Get product
PUT    /api/v1/products/:id          # Update product
DELETE /api/v1/products/:id          # Delete product
GET    /api/v1/products              # List with pagination
POST   /api/v1/products/bulk-update-status
GET    /api/v1/products/stats        # Statistics
```

## Working with Generated Projects

### Start Development

```bash
cd your-project

# Start all services (MySQL, Redis, Kafka)
make up

# Run database migrations
make migrate-up

# Start the application
make run
```

### Available Make Commands

```bash
make up              # Start Docker services
make down            # Stop Docker services
make migrate-up      # Run migrations
make migrate-down    # Rollback migrations
make run             # Run application
make build           # Build binary
make test            # Run tests
make lint            # Run linter
make clean           # Clean artifacts
make logs            # View service logs
```

### Configuration

Projects support 3 configuration layers:

1. **Environment variables** (highest priority)
   ```bash
   export DB_HOST=localhost
   export DB_PORT=3306
   export APP_PORT=8080
   ```

2. **config.toml file**
   ```toml
   [database]
   host = "localhost"
   port = 3306
   
   [server]
   port = 8080
   ```

3. **Default values** (fallback)

### Project Structure

```
your-project/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── application.go          # App setup
│   ├── config/                 # Configuration
│   ├── database/               # DB + migrations
│   ├── entity/                 # Domain models
│   ├── dto/                    # Request/Response
│   ├── repository/             # Data access
│   ├── service/                # Business logic
│   └── handler/                # HTTP handlers
├── docker-compose.yml          # Services
├── Dockerfile                  # App container
├── Makefile                    # Commands
└── config.toml                 # Config file
```

## Examples

### E-commerce Product API
```bash
ready-go new product-service \
  --module github.com/myshop/products \
  --sample-api Product \
  --author "Shop Team"
```

### User Management Service
```bash
ready-go new user-service \
  --module github.com/mycompany/users \
  --sample-api User
```

### Lightweight Microservice (No Kafka/Redis)
```bash
ready-go new orders-api \
  --module github.com/store/orders \
  --sample-api Order \
  --with-kafka=false \
  --with-redis=false
```

### Interactive Setup
```bash
ready-go new my-project -i
# Follow the prompts to configure your project
```

## Testing Your API

Once your app is running, test it:

```bash
# Health check
curl http://localhost:8080/health

# Create a record (e.g., Product)
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"Gaming laptop","status":"active"}'

# List records
curl http://localhost:8080/api/v1/products

# Get stats
curl http://localhost:8080/api/v1/products/stats
```

## Troubleshooting

### Port Already in Use
```bash
# Check what's using the port
lsof -i :8080

# Or change port in config.toml
[server]
port = 9090
```

### Database Connection Failed
```bash
# Make sure MySQL is running
docker ps | grep mysql

# Check connection in config.toml or .env.docker
```

### Dependencies Not Downloaded
```bash
cd your-project
go mod tidy
```

## Tips

1. **Customize Entity**: Use `--sample-api` to set your domain model name
2. **Module Path**: Always use your GitHub/GitLab username in `--module`
3. **Git Init**: Projects auto-initialize git (use `--skip-git` to disable)
4. **Docker First**: Always run `make up` before `make run`
5. **Migrations**: Run `make migrate-up` after starting Docker services

## Getting Help

```bash
# Show help
ready-go --help

# Show version
ready-go --version

# Command help
ready-go new --help
```

## Next Steps

After generating a project:

1. **Review** the generated README.md in your project
2. **Customize** entity fields in `internal/entity/`
3. **Add** more endpoints in handlers
4. **Write** tests for your business logic
5. **Deploy** using the included Dockerfile

## Learn More

- All projects follow **Clean Architecture** principles
- Configuration uses **multi-source loading** (ENV → TOML → Defaults)
- API uses **Fiber v2** framework
- Database uses **sqlx** with MySQL
- Includes **health checks** for Kubernetes/monitoring

---

**Happy coding!** 🚀

For issues or questions, check the project README.md or visit the repository.
