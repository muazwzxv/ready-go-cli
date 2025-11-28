package generator

import (
	"embed"
	"fmt"
	"os"
	"text/template"

	"github.com/muazwzxv/ready-go-cli/internal/config"
)

var embeddedFS embed.FS

// SetEmbeddedTemplates sets the embedded filesystem for templates
func SetEmbeddedTemplates(fs embed.FS) {
	embeddedFS = fs
}

// TemplateRenderer handles rendering of template files
type TemplateRenderer struct {
	config *config.ProjectConfig
}

// NewTemplateRenderer creates a new TemplateRenderer
func NewTemplateRenderer(cfg *config.ProjectConfig) *TemplateRenderer {
	return &TemplateRenderer{
		config: cfg,
	}
}

// RenderToFile renders a template to a file
func (r *TemplateRenderer) RenderToFile(templateName, outputPath string) error {
	// Read template content
	content, err := r.getTemplateContent(templateName)
	if err != nil {
		return fmt.Errorf("failed to get template %s: %w", templateName, err)
	}

	// Parse and execute template
	tmpl, err := template.New(templateName).Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, r.config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// getTemplateContent retrieves template content from embedded filesystem
func (r *TemplateRenderer) getTemplateContent(templateName string) (string, error) {
	// Try to read from embedded filesystem first
	path := "templates/" + templateName
	data, err := embeddedFS.ReadFile(path)
	if err == nil {
		return string(data), nil
	}

	// Fallback to reading from disk (for development)
	diskPaths := []string{
		"templates/" + templateName,
		"../../templates/" + templateName,
		templateName,
	}

	for _, diskPath := range diskPaths {
		data, err := os.ReadFile(diskPath)
		if err == nil {
			return string(data), nil
		}
	}

	return "", fmt.Errorf("template %s not found in embedded FS or disk", templateName)
}
