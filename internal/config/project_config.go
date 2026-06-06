package config

import (
	"fmt"
	"regexp"
	"strings"
)

// ProjectConfig holds all configuration for generating a new project
type ProjectConfig struct {
	ProjectName        string
	ModuleName         string
	GoVersion          string
	OutputDir          string
	Router             string
	SampleAPIName      string
	SampleAPINameLower string
	SampleAPINameUpper string
	SampleTableName    string
	ServerPort         string
	DBPort             string
	RedisPort          string
	KafkaPort          string
}

// NewProjectConfig creates a new ProjectConfig with default values
func NewProjectConfig(projectName string) *ProjectConfig {
	return &ProjectConfig{
		ProjectName:   projectName,
		ModuleName:    fmt.Sprintf("github.com/username/%s", projectName),
		GoVersion:     "1.21",
		OutputDir:     ".",
		Router:        "fiber",
		SampleAPIName: "User",
		ServerPort:    "8080",
		DBPort:        "3306",
		RedisPort:     "6379",
		KafkaPort:     "9092",
	}
}

// TemplateDir returns the template directory prefix for the chosen router
func (c *ProjectConfig) TemplateDir() string {
	return c.Router
}

// Validate checks if the configuration is valid
func (c *ProjectConfig) Validate() error {
	if c.ProjectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, c.ProjectName)
	if !matched {
		return fmt.Errorf("project name can only contain letters, numbers, hyphens, and underscores")
	}

	if c.ModuleName == "" {
		return fmt.Errorf("module name cannot be empty")
	}

	if c.SampleAPIName == "" {
		return fmt.Errorf("sample API name cannot be empty")
	}

	if c.Router != "fiber" && c.Router != "chi" {
		return fmt.Errorf("router must be either 'fiber' or 'chi'")
	}

	return nil
}

// Process calculates derived fields from the configuration
func (c *ProjectConfig) Process() {
	c.SampleAPIName = strings.Title(c.SampleAPIName)
	c.SampleAPINameLower = strings.ToLower(c.SampleAPIName)
	c.SampleAPINameUpper = strings.ToUpper(c.SampleAPIName)
	c.SampleTableName = pluralize(c.SampleAPINameLower)
}

func pluralize(word string) string {
	if strings.HasSuffix(word, "y") {
		return strings.TrimSuffix(word, "y") + "ies"
	}
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
		return word + "es"
	}
	return word + "s"
}
