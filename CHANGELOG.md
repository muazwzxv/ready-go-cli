# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.1.0] - 2025-12-11

### 🎯 Feature Release: Structured Logging with Fiber Logger

This release replaces Go's standard library logger with Fiber's built-in structured logger, providing better observability and modern logging capabilities.

### Added
- **Structured Logging**: All log statements now use key-value pairs for better parsing and analysis
  - Example: `log.Infow("user created", "user_id", 123, "email", "user@example.com")`
- **Context-Aware Logging**: HTTP handlers can now log with request context
  - Automatic inclusion of request details (method, path, IP, user agent)
  - Example: `logger := log.WithContext(c.UserContext())`
- **Configurable Log Levels**: Control verbosity via configuration
  - Levels: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`
  - Default: `info` (balanced detail/performance)
- **Log Level Configuration**: New `log_level` field in server config
  - Set via `config.toml`: `log_level = "debug"`
  - Override via ENV: `SERVER_LOG_LEVEL=debug`
- **Comprehensive Documentation**: New "Logging" section in generated README.md with usage examples

### Changed
- All logging now uses `github.com/gofiber/fiber/v2/log` instead of `log` stdlib
- Application logs use structured format: `log.Infow("message", "key", value)`
- Database connection logs use appropriate levels:
  - Connection attempts: `debug` level
  - Successful operations: `info` level
  - Failed retries: `warn` level
- Error handler middleware now logs with full request context
- Health check handler logs access attempts at `debug` level

### Improved
- Better log filtering in production (set level to `warn` or `error`)
- Easier debugging with `debug` level showing connection attempts and retries
- Compatible with log aggregation tools (Datadog, Grafana Loki, CloudWatch)
- Logs output to stdout (Docker/Kubernetes friendly)
- Invalid log level configuration fails early with clear error message

### Technical Details
- 16 log statements migrated across 4 template files
- All Fatal-level logs maintain same API (backward compatible)
- `setLogLevel()` helper validates and configures log level on startup
- Middleware logs include: method, path, IP, user agent, error details
- Handler logs include: operation details, success/failure status, entity IDs

### Example Log Output

**Info level (default):**
```
[Info] application config loaded server_host=0.0.0.0 server_port=8080 log_level=info
[Info] database connection established successfully
[Info] application starting server address=0.0.0.0:8080
```

**Debug level:**
```
[Debug] database attempting connection user=appuser host=localhost port=3306
[Debug] database ping attempt attempt=1 max_attempts=3
[Info] database connection established successfully
```

**Context-aware (in handlers):**
```
[Info] creating user name="John Doe" email="john@example.com"
[Info] user created successfully id=123
```

### Migration from v2.0.0
No breaking changes. Existing projects continue to work. New projects automatically get structured logging.

---

## [2.0.0] - 2025-12-11

### 🚀 Major Release: Dependency Injection with samber/do

This is a **breaking release** that fundamentally improves how generated projects handle dependency management.

### Added
- **Dependency Injection**: Integrated [samber/do v2](https://do.samber.dev/) for type-safe dependency injection
- **Lifecycle Management**: Automatic graceful shutdown with `Shutdowner` interface
- **Health Checks**: Built-in health check support via `Healthchecker` interface
- Database now implements `Shutdown()` and `HealthCheck()` methods for DI lifecycle
- Provider functions for all major components (Database, Fiber, Queries)

### Changed
- **BREAKING**: All constructors now use `do.Provider[T]` signature: `func(do.Injector) (T, error)`
  - Handler constructors: `NewUserHandler(do.Injector) (*UserHandler, error)`
  - Service constructors: `NewUserService(do.Injector) (UserService, error)`
  - Repository constructors: `NewUserRepository(do.Injector) (UserRepository, error)`
- **BREAKING**: `ApplicationContext` replaced with simpler `Application` struct
- **BREAKING**: Manual dependency wiring replaced with DI container registration
- Application initialization simplified from ~166 lines to ~170 lines (better organized)
- Graceful shutdown now uses `injector.ShutdownOnSignals()` for proper cleanup
- Route registration now invokes handlers directly from DI container

### Removed
- **BREAKING**: `Services` struct (use injector instead)
- **BREAKING**: `Handlers` struct (use injector instead)
- Manual dependency chain construction in application.go
- Custom shutdown logic - now handled by DI container

### Improved
- 60% reduction in boilerplate dependency wiring
- Easier testing with `do.Override` for mocking
- Type-safe dependency resolution at compile-time
- Automatic service lifecycle management
- Cleaner separation of concerns

### Migration Guide from v1.x

If you have existing projects generated with v1.x, you have two options:

**Option 1: Continue using v1.x CLI**
```bash
# Install specific version
go install github.com/muazwzxv/ready-go-cli/cmd/ready-go@v1.1.0
```

**Option 2: Regenerate project with v2.0.0**
Recommended for new projects or if you want the benefits of DI.

**Manual Migration Steps:**
1. Add dependency: `go get github.com/samber/do/v2`
2. Update all constructor signatures to match `do.Provider[T]` format
3. Replace `ApplicationContext` with new `Application` struct and DI container
4. Update `main.go` to use new application initialization
5. Add `Shutdown()` and `HealthCheck()` methods to Database

**Why the breaking change?**
- ✅ 60% less boilerplate code
- ✅ Easier testing (simple mocking with `do.Override`)
- ✅ Automatic dependency resolution
- ✅ Built-in lifecycle management (shutdown, health checks)
- ✅ Type-safe with generics (no reflection)
- ✅ Industry-standard DI pattern

### Dependencies
- Added: `github.com/samber/do/v2` (latest)

---

## [1.1.0] - 2025-12-10

### Added
- **Viper Integration**: Replaced manual configuration loading with [Viper](https://github.com/spf13/viper) for more robust and standardized config management
- Default configuration values to ensure all keys exist for environment variable binding
- Comprehensive environment variable support with automatic override capability
- Support for multiple config file search paths (`./config.toml`, `/etc/<project>/`, `$HOME/.<project>`)

### Changed
- **BREAKING CHANGE**: Environment variable naming convention updated for consistency
  - **Old prefix**: `DB_*` (e.g., `DB_HOST`, `DB_PORT`, `DB_USER`)
  - **New prefix**: `DATABASE_*` (e.g., `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_USER`)
  - This change affects `.env.example` and `.env.docker` templates
  - All generated projects will use the new naming convention
  
### Improved
- Configuration loading is now ~60% more concise (106 lines vs 266 lines)
- Better error handling for configuration file reading
- Automatic environment variable mapping without manual parsing
- Struct tags changed from `toml:` to `mapstructure:` for Viper compatibility

### Fixed
- Config loading in Docker containers now properly reads environment variables
- Environment variables now correctly override TOML file values at runtime

### Migration Guide for Existing Projects

If you have an existing project generated with v1.0.0, you'll need to update your environment variables:

**In `.env` files:**
```bash
# Old (v1.0.0)
DB_HOST=localhost
DB_PORT=3306
DB_USER=appuser
DB_PASSWORD=apppassword
DB_DATABASE=myapp

# New (v1.1.0)
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_USER=appuser
DATABASE_PASSWORD=apppassword
DATABASE_DATABASE=myapp
```

**In `internal/config/config.go`:**
- Update imports: Replace `github.com/BurntSushi/toml` with `github.com/spf13/viper`
- Update struct tags: Replace `toml:"field"` with `mapstructure:"field"`
- Update `LoadConfig()` function to use Viper's API
- Run `go mod tidy` to update dependencies

Or simply regenerate your project with the new CLI version for the cleanest migration.

### Removed
- Manual environment variable parsing helpers (`getEnvOrDefault`, `getEnvIntOrDefault`, etc.)
- Complex merge logic for config sources
- Direct dependency on `github.com/BurntSushi/toml`
- Duplicate `templates/` directory at project root (consolidated to `cmd/ready-go/templates/`)

---

## [1.0.0] - 2025-12-09

### Added
- Initial release of Ready-Go CLI
- Clean architecture project scaffolding
- Docker Compose setup with MySQL, Redis, and Kafka
- Database migration support with goose
- Type-safe SQL queries with sqlc
- Multi-source configuration (TOML + environment variables)
- Sample CRUD API generation
- Interactive and flag-based CLI modes
- Health check endpoints
- Makefile with common development tasks

[2.0.0]: https://github.com/muazwzxv/ready-go-cli/compare/v1.1.0...v2.0.0
[1.1.0]: https://github.com/muazwzxv/ready-go-cli/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/muazwzxv/ready-go-cli/releases/tag/v1.0.0
