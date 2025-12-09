# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[1.1.0]: https://github.com/muazwzxv/ready-go-cli/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/muazwzxv/ready-go-cli/releases/tag/v1.0.0
