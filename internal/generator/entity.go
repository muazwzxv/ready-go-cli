package generator

import (
	"fmt"
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

// Generate creates the entity file and migration
func (g *EntityGenerator) Generate() error {
	// Generate entity file
	if err := g.generateEntityFile(); err != nil {
		return fmt.Errorf("generate entity file: %w", err)
	}

	// Generate migration file
	if err := g.generateMigrationFile(); err != nil {
		return fmt.Errorf("generate migration file: %w", err)
	}

	return nil
}

// generateEntityFile creates the entity Go file
func (g *EntityGenerator) generateEntityFile() error {
	outputPath := filepath.Join(
		g.config.ProjectPath,
		"internal",
		"entity",
		g.config.EntityNameLower+".go",
	)

	err := RenderTemplateToFile(
		"internal/entity/entity.go.tmpl",
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
	// Generate timestamp in goose format (YYYYMMDDHHMMSS)
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_create_%s.sql", timestamp, g.config.TableName)

	outputPath := filepath.Join(
		g.config.ProjectPath,
		"internal",
		"database",
		"migrations",
		filename,
	)

	err := RenderTemplateToFile(
		"internal/database/migrations/migration.sql.tmpl",
		outputPath,
		g.config.TemplateData(),
	)
	if err != nil {
		return err
	}

	fmt.Printf("  ✓ Created %s\n", outputPath)
	return nil
}
