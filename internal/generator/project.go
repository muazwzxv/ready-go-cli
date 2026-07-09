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
	fmt.Println("🔀 Initializing git repository...")
	if err := g.initGit(projectPath); err != nil {
		fmt.Printf("⚠️  Warning: failed to initialize git: %v\n", err)
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
		filepath.Join(projectPath, "cmd", "api"),
		filepath.Join(projectPath, "internal", "config"),
		filepath.Join(projectPath, "internal", "handlers", g.config.SampleAPINameLower),
		filepath.Join(projectPath, "internal", "handlers", "util"),
		filepath.Join(projectPath, "internal", "models"),
		filepath.Join(projectPath, "internal", "repository"),
		filepath.Join(projectPath, "internal", "telemetry"),
		filepath.Join(projectPath, "database", "migrations"),
		filepath.Join(projectPath, "database", "queries"),
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

	// Router-specific template prefix
	routerPrefix := g.config.TemplateDir()

	files := []struct {
		template string
		output   string
	}{
		// Router-specific Go source files
		{filepath.Join(routerPrefix, "cmd/api/main.go.tmpl"), filepath.Join(projectPath, "cmd", "api", "main.go")},
		{filepath.Join(routerPrefix, "internal/handlers/handler.go.tmpl"), filepath.Join(projectPath, "internal", "handlers", "handler.go")},
		{filepath.Join(routerPrefix, "internal/handlers/sample/sample_handler.go.tmpl"), filepath.Join(projectPath, "internal", "handlers", g.config.SampleAPINameLower, "handler.go")},
		{filepath.Join(routerPrefix, "internal/handlers/util/util.go.tmpl"), filepath.Join(projectPath, "internal", "handlers", "util", "util.go")},

		// Shared Go source files
		{"shared/cmd/service.go.tmpl", filepath.Join(projectPath, "cmd", "service.go")},
		{"shared/internal/config/config.go.tmpl", filepath.Join(projectPath, "internal", "config", "config.go")},
		{"shared/internal/models/db.go.tmpl", filepath.Join(projectPath, "internal", "models", "db.go")},
		{"shared/internal/repository/db.go.tmpl", filepath.Join(projectPath, "internal", "repository", "db.go")},
		{"shared/internal/telemetry/telemetry.go.tmpl", filepath.Join(projectPath, "internal", "telemetry", "telemetry.go")},

		// Database files
		{"shared/database/migrations/init.sql.tmpl", filepath.Join(projectPath, "database", "migrations", "00001_init.sql")},
		{"shared/database/queries/sample.sql.tmpl", filepath.Join(projectPath, "database", "queries", g.config.SampleAPINameLower+".sql")},

		// Project config files
		{"shared/project/docker-compose.yml.tmpl", filepath.Join(projectPath, "docker-compose.yml")},
		{"shared/project/Dockerfile.tmpl", filepath.Join(projectPath, "Dockerfile")},
		{"shared/project/Makefile.tmpl", filepath.Join(projectPath, "Makefile")},
		{"shared/project/sqlc.yaml.tmpl", filepath.Join(projectPath, "sqlc.yaml")},
		{"shared/project/otel-collector-config.yaml", filepath.Join(projectPath, "otel-collector-config.yaml")},
		{"shared/project/.env.example.tmpl", filepath.Join(projectPath, ".env.example")},
		{"shared/project/README.md.tmpl", filepath.Join(projectPath, "README.md")},
	}

	// Chi-specific extra files
	if g.config.Router == "chi" {
		files = append(files, struct {
			template string
			output   string
		}{
			template: filepath.Join(routerPrefix, "internal/handlers/util/json.go.tmpl"),
			output:   filepath.Join(projectPath, "internal", "handlers", "util", "json.go"),
		})
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
