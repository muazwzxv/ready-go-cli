package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// EntityConfig holds configuration for generating a new entity
type EntityConfig struct {
	EntityName      string // PascalCase: "Product"
	EntityNameLower string // lowercase: "product"
	TableName       string // plural: "products"
	ProjectPath     string // current working directory
}

// NewEntityConfig creates a new EntityConfig with the given entity name
func NewEntityConfig(entityName string) *EntityConfig {
	return &EntityConfig{
		EntityName: entityName,
	}
}

// Process calculates derived fields from the configuration
func (c *EntityConfig) Process() {
	// Ensure first letter is uppercase (PascalCase)
	if len(c.EntityName) > 0 {
		c.EntityName = strings.ToUpper(c.EntityName[:1]) + c.EntityName[1:]
	}
	c.EntityNameLower = strings.ToLower(c.EntityName)
	c.TableName = pluralize(c.EntityNameLower)

	if c.ProjectPath == "" {
		c.ProjectPath, _ = os.Getwd()
	}
}

// Validate checks if the configuration is valid and project structure exists
func (c *EntityConfig) Validate() error {
	// Check entity name is provided
	if c.EntityName == "" {
		return fmt.Errorf("entity name cannot be empty")
	}

	// Check PascalCase (starts with uppercase, alphanumeric only)
	if !regexp.MustCompile(`^[A-Z][a-zA-Z0-9]*$`).MatchString(c.EntityName) {
		return fmt.Errorf("entity name must be PascalCase (e.g., Product, OrderItem)")
	}

	// Check project structure exists
	if _, err := os.Stat(filepath.Join(c.ProjectPath, "go.mod")); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found - run this command from a Go project root")
	}

	// Check database directories exist (simplified structure)
	migrationsDir := filepath.Join(c.ProjectPath, "database", "migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("database/migrations directory not found - is this a ready-go project?")
	}

	queriesDir := filepath.Join(c.ProjectPath, "database", "queries")
	if _, err := os.Stat(queriesDir); os.IsNotExist(err) {
		return fmt.Errorf("database/queries directory not found - is this a ready-go project?")
	}

	// Check entity file doesn't already exist
	entityFile := filepath.Join(c.ProjectPath, "internal", "entity", c.EntityNameLower+".go")
	if _, err := os.Stat(entityFile); err == nil {
		return fmt.Errorf("entity %s already exists at %s", c.EntityName, entityFile)
	}

	return nil
}

// TemplateData returns a map compatible with existing templates
func (c *EntityConfig) TemplateData() map[string]any {
	return map[string]any{
		"EntityName":      c.EntityName,
		"EntityNameLower": c.EntityNameLower,
		"TableName":       c.TableName,
	}
}
