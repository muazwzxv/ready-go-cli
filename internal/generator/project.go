package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/muazwzxv/ready-go-cli/internal/config"
)

// ProjectGenerator handles the generation of a new project
type ProjectGenerator struct {
	config *config.ProjectConfig
}

// NewProjectGenerator creates a new ProjectGenerator
func NewProjectGenerator(cfg *config.ProjectConfig) *ProjectGenerator {
	return &ProjectGenerator{
		config: cfg,
	}
}

// Generate creates the entire project structure
func (g *ProjectGenerator) Generate() error {
	projectPath := filepath.Join(g.config.OutputDir, g.config.ProjectName)

	// Check if project directory already exists
	if _, err := os.Stat(projectPath); err == nil {
		return fmt.Errorf("directory %s already exists", projectPath)
	}

	// Create directory structure
	fmt.Println("📁 Creating directory structure...")
	if err := g.createDirectoryStructure(projectPath); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Generate files
	fmt.Println("📝 Generating project files...")
	if err := g.generateFiles(projectPath); err != nil {
		return fmt.Errorf("failed to generate files: %w", err)
	}

	// Initialize go module
	fmt.Println("🔧 Initializing go module...")
	if err := g.initGoModule(projectPath); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	// Initialize git repository
	if !g.config.SkipGit {
		fmt.Println("🔀 Initializing git repository...")
		if err := g.initGit(projectPath); err != nil {
			fmt.Printf("⚠️  Warning: failed to initialize git: %v\n", err)
		}
	}

	// Download dependencies
	fmt.Println("📦 Downloading dependencies...")
	if err := g.downloadDependencies(projectPath); err != nil {
		fmt.Printf("⚠️  Warning: failed to download dependencies: %v\n", err)
		fmt.Println("   Run 'go mod tidy' manually in the project directory")
	}

	return nil
}

// createDirectoryStructure creates the project directory structure
func (g *ProjectGenerator) createDirectoryStructure(projectPath string) error {
	dirs := []string{
		projectPath,
		filepath.Join(projectPath, "cmd", "server"),
		filepath.Join(projectPath, "internal", "config"),
		filepath.Join(projectPath, "internal", "database", "migrations"),
		filepath.Join(projectPath, "internal", "database", "query"),
		filepath.Join(projectPath, "internal", "database", "store"),
		filepath.Join(projectPath, "internal", "entity"),
		filepath.Join(projectPath, "internal", "dto", "request"),
		filepath.Join(projectPath, "internal", "dto", "response"),
		filepath.Join(projectPath, "internal", "repository"),
		filepath.Join(projectPath, "internal", "service", g.config.SampleAPINameLower),
		filepath.Join(projectPath, "internal", "handler", "health"),
		filepath.Join(projectPath, "internal", "handler", g.config.SampleAPINameLower),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateFiles generates all project files from templates
func (g *ProjectGenerator) generateFiles(projectPath string) error {
	renderer := NewTemplateRenderer(g.config)

	files := []struct {
		template string
		output   string
	}{
		// Root level files
		{"project/main.go.tmpl", filepath.Join(projectPath, "cmd", "server", "main.go")},
		{"project/Dockerfile.tmpl", filepath.Join(projectPath, "Dockerfile")},
		{"project/docker-compose.yml.tmpl", filepath.Join(projectPath, "docker-compose.yml")},
		{"project/Makefile.tmpl", filepath.Join(projectPath, "Makefile")},
		{"project/config.toml.tmpl", filepath.Join(projectPath, "config.toml")},
		{"project/.env.docker.tmpl", filepath.Join(projectPath, ".env.docker")},
		{"project/.env.example.tmpl", filepath.Join(projectPath, ".env.example")},
		{"project/.gitignore.tmpl", filepath.Join(projectPath, ".gitignore")},
		{"project/README.md.tmpl", filepath.Join(projectPath, "README.md")},

		// Internal files
		{"internal/application.go.tmpl", filepath.Join(projectPath, "internal", "application.go")},
		{"internal/config/config.go.tmpl", filepath.Join(projectPath, "internal", "config", "config.go")},
		{"internal/database/database.go.tmpl", filepath.Join(projectPath, "internal", "database", "database.go")},
		{"internal/database/sqlc.yaml.tmpl", filepath.Join(projectPath, "internal", "database", "sqlc.yaml")},
		{"internal/database/migrations/migration.sql.tmpl", filepath.Join(projectPath, "internal", "database", "migrations", "001_create_initial_schema.sql")},
		{"internal/database/query/sample.sql.tmpl", filepath.Join(projectPath, "internal", "database", "query", g.config.SampleAPINameLower+".sql")},
		{"internal/database/store/.gitignore.tmpl", filepath.Join(projectPath, "internal", "database", "store", ".gitignore")},

		// Entity
		{"internal/entity/entity.go.tmpl", filepath.Join(projectPath, "internal", "entity", g.config.SampleAPINameLower+".go")},

		// DTO
		{"internal/dto/request/common.go.tmpl", filepath.Join(projectPath, "internal", "dto", "request", "common.go")},
		{"internal/dto/request/sample_request.go.tmpl", filepath.Join(projectPath, "internal", "dto", "request", g.config.SampleAPINameLower+"_request.go")},
		{"internal/dto/response/common.go.tmpl", filepath.Join(projectPath, "internal", "dto", "response", "common.go")},
		{"internal/dto/response/error_response.go.tmpl", filepath.Join(projectPath, "internal", "dto", "response", "error_response.go")},
		{"internal/dto/response/sample_response.go.tmpl", filepath.Join(projectPath, "internal", "dto", "response", g.config.SampleAPINameLower+"_response.go")},

		// Repository
		{"internal/repository/interfaces.go.tmpl", filepath.Join(projectPath, "internal", "repository", "interfaces.go")},
		{"internal/repository/sample_repository.go.tmpl", filepath.Join(projectPath, "internal", "repository", g.config.SampleAPINameLower+"_repository.go")},

		// Service
		{"internal/service/sample/sample.go.tmpl", filepath.Join(projectPath, "internal", "service", g.config.SampleAPINameLower, g.config.SampleAPINameLower+".go")},
		{"internal/service/sample/create_sample_service.go.tmpl", filepath.Join(projectPath, "internal", "service", g.config.SampleAPINameLower, "create_"+g.config.SampleAPINameLower+"_service.go")},

		// Handler
		{"internal/handler/health/health_handler.go.tmpl", filepath.Join(projectPath, "internal", "handler", "health", "health_handler.go")},
		{"internal/handler/middleware.go.tmpl", filepath.Join(projectPath, "internal", "handler", "middleware.go")},
		{"internal/handler/sample/sample_handler.go.tmpl", filepath.Join(projectPath, "internal", "handler", g.config.SampleAPINameLower, g.config.SampleAPINameLower+"_handler.go")},
		{"internal/handler/sample/create_sample_handler.go.tmpl", filepath.Join(projectPath, "internal", "handler", g.config.SampleAPINameLower, "create_"+g.config.SampleAPINameLower+"_handler.go")},
	}

	for _, file := range files {
		if err := renderer.RenderToFile(file.template, file.output); err != nil {
			return fmt.Errorf("failed to generate %s: %w", file.output, err)
		}
	}

	return nil
}

// initGoModule initializes the go module
func (g *ProjectGenerator) initGoModule(projectPath string) error {
	cmd := exec.Command("go", "mod", "init", g.config.ModuleName)
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod init failed: %w\n%s", err, output)
	}
	return nil
}

// downloadDependencies downloads all project dependencies
func (g *ProjectGenerator) downloadDependencies(projectPath string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w\n%s", err, output)
	}
	return nil
}

// initGit initializes a git repository
func (g *ProjectGenerator) initGit(projectPath string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git init failed: %w\n%s", err, output)
	}
	return nil
}
