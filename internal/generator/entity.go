package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/muazwzxv/ready-go-cli/internal/config"
)

// EntityGenerator handles generation of new entities in existing projects
type EntityGenerator struct {
	config *config.EntityConfig
}

// NewEntityGenerator creates a new EntityGenerator
func NewEntityGenerator(cfg *config.EntityConfig) *EntityGenerator {
	return &EntityGenerator{
		config: cfg,
	}
}

// Generate creates the entity file, migration, and queries
func (g *EntityGenerator) Generate() error {
	// Generate entity file
	if err := g.generateEntityFile(); err != nil {
		return fmt.Errorf("generate entity file: %w", err)
	}

	// Generate migration file
	if err := g.generateMigrationFile(); err != nil {
		return fmt.Errorf("generate migration file: %w", err)
	}

	// Generate queries file
	if err := g.generateQueriesFile(); err != nil {
		return fmt.Errorf("generate queries file: %w", err)
	}

	return nil
}

// generateEntityFile creates the entity Go file
func (g *EntityGenerator) generateEntityFile() error {
	// Create internal/entity directory if it doesn't exist
	entityDir := filepath.Join(g.config.ProjectPath, "internal", "entity")
	if err := createDirIfNotExists(entityDir); err != nil {
		return err
	}

	outputPath := filepath.Join(entityDir, g.config.EntityNameLower+".go")

	err := RenderTemplateToFile(
		"entity/entity.go.tmpl",
		outputPath,
		g.config.TemplateData(),
	)
	if err != nil {
		return err
	}

	fmt.Printf("  ✓ Created %s\n", outputPath)
	return nil
}

// generateMigrationFile creates the goose migration SQL file
func (g *EntityGenerator) generateMigrationFile() error {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_create_%s.sql", timestamp, g.config.TableName)

	outputPath := filepath.Join(
		g.config.ProjectPath,
		"database",
		"migrations",
		filename,
	)

	err := RenderTemplateToFile(
		"entity/migration.sql.tmpl",
		outputPath,
		g.config.TemplateData(),
	)
	if err != nil {
		return err
	}

	fmt.Printf("  ✓ Created %s\n", outputPath)
	return nil
}

// generateQueriesFile creates the SQLC queries file
func (g *EntityGenerator) generateQueriesFile() error {
	outputPath := filepath.Join(
		g.config.ProjectPath,
		"database",
		"queries",
		g.config.EntityNameLower+".sql",
	)

	err := RenderTemplateToFile(
		"entity/queries.sql.tmpl",
		outputPath,
		g.config.TemplateData(),
	)
	if err != nil {
		return err
	}

	fmt.Printf("  ✓ Created %s\n", outputPath)
	return nil
}

func createDirIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
