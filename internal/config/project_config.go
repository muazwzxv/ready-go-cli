package config

import (
	"fmt"
	"regexp"
	"strings"
)

// ProjectConfig holds all configuration for generating a new project
type ProjectConfig struct {
	// Project metadata
	ProjectName  string
	ModuleName   string
	Description  string
	Author       string
	GoVersion    string
	
	// Output configuration
	OutputDir    string
	
	// Service options
	WithRedis    bool
	WithKafka    bool
	
	// Sample API configuration
	SampleAPIName      string  // e.g., "User", "Product"
	SampleAPINameLower string  // e.g., "user", "product"
	SampleTableName    string  // e.g., "users", "products"
	
	// Docker configuration
	MySQLVersion  string
	RedisVersion  string
	KafkaVersion  string
	
	// Port configuration
	AppPort       int
	MySQLPort     int
	RedisPort     int
	KafkaPort     int
	KafkaUIPort   int
	
	// Git configuration
	SkipGit       bool
}

// NewProjectConfig creates a new ProjectConfig with default values
func NewProjectConfig(projectName string) *ProjectConfig {
	return &ProjectConfig{
		ProjectName:   projectName,
		ModuleName:    fmt.Sprintf("github.com/username/%s", projectName),
		Description:   fmt.Sprintf("A %s service", projectName),
		Author:        "",
		GoVersion:     "1.23",
		OutputDir:     ".",
		WithRedis:     true,
		WithKafka:     true,
		SampleAPIName: "User",
		MySQLVersion:  "8.0",
		RedisVersion:  "7-alpine",
		KafkaVersion:  "7.6.0",
		AppPort:       8080,
		MySQLPort:     3306,
		RedisPort:     6379,
		KafkaPort:     9092,
		KafkaUIPort:   8090,
		SkipGit:       false,
	}
}

// Validate checks if the configuration is valid
func (c *ProjectConfig) Validate() error {
	if c.ProjectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	
	// Validate project name (alphanumeric, hyphens, underscores)
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
	
	return nil
}

// Process calculates derived fields from the configuration
func (c *ProjectConfig) Process() {
	// Process sample API name
	c.SampleAPIName = strings.Title(c.SampleAPIName)
	c.SampleAPINameLower = strings.ToLower(c.SampleAPIName)
	
	// Generate table name (pluralize)
	c.SampleTableName = pluralize(c.SampleAPINameLower)
}

// pluralize converts a singular word to plural (simple implementation)
func pluralize(word string) string {
	if strings.HasSuffix(word, "y") {
		return strings.TrimSuffix(word, "y") + "ies"
	}
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
		return word + "es"
	}
	return word + "s"
}

// GetContainerName generates a container name for Docker services
func (c *ProjectConfig) GetContainerName(service string) string {
	return fmt.Sprintf("%s-%s", c.ProjectName, service)
}

// GetNetworkName generates a Docker network name
func (c *ProjectConfig) GetNetworkName() string {
	return fmt.Sprintf("%s-network", c.ProjectName)
}

// GetVolumeName generates a Docker volume name
func (c *ProjectConfig) GetVolumeName(service string) string {
	return fmt.Sprintf("%s-%s-data", c.ProjectName, service)
}
