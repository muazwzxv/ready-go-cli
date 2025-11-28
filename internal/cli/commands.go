package cli

import (
	"fmt"

	"github.com/muazwzxv/ready-go-cli/internal/config"
	"github.com/muazwzxv/ready-go-cli/internal/generator"
	"github.com/urfave/cli/v2"
)

// Commands returns all CLI commands
func Commands() []*cli.Command {
	return []*cli.Command{
		NewCommand(),
	}
}

// NewCommand creates the 'new' command for scaffolding projects
func NewCommand() *cli.Command {
	return &cli.Command{
		Name:      "new",
		Usage:     "Create a new Go project",
		ArgsUsage: "<project-name>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "module",
				Aliases: []string{"m"},
				Usage:   "Go module name",
			},
			&cli.StringFlag{
				Name:    "description",
				Aliases: []string{"d"},
				Usage:   "Project description",
			},
			&cli.StringFlag{
				Name:  "author",
				Usage: "Author name",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   ".",
				Usage:   "Output directory",
			},
			&cli.BoolFlag{
				Name:  "with-redis",
				Value: true,
				Usage: "Include Redis in docker-compose",
			},
			&cli.BoolFlag{
				Name:  "with-kafka",
				Value: true,
				Usage: "Include Kafka in docker-compose",
			},
			&cli.StringFlag{
				Name:  "sample-api",
				Value: "User",
				Usage: "Sample API entity name (e.g., User, Product)",
			},
			&cli.BoolFlag{
				Name:  "skip-git",
				Usage: "Skip git initialization",
			},
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Usage:   "Interactive mode with prompts",
			},
		},
		Action: newProjectAction,
	}
}

// newProjectAction handles the 'new' command execution
func newProjectAction(c *cli.Context) error {
	// Get project name from arguments
	projectName := c.Args().First()
	
	if projectName == "" {
		return fmt.Errorf("project name is required\nUsage: ready-go new <project-name>")
	}
	
	// Create project configuration
	cfg := config.NewProjectConfig(projectName)
	
	// Check if interactive mode
	if c.Bool("interactive") {
		if err := PromptForConfig(cfg); err != nil {
			return fmt.Errorf("interactive prompt failed: %w", err)
		}
	} else {
		// Apply flags
		if c.String("module") != "" {
			cfg.ModuleName = c.String("module")
		}
		if c.String("description") != "" {
			cfg.Description = c.String("description")
		}
		if c.String("author") != "" {
			cfg.Author = c.String("author")
		}
		cfg.OutputDir = c.String("output")
		cfg.WithRedis = c.Bool("with-redis")
		cfg.WithKafka = c.Bool("with-kafka")
		cfg.SampleAPIName = c.String("sample-api")
		cfg.SkipGit = c.Bool("skip-git")
	}
	
	// Process configuration
	cfg.Process()
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	
	// Generate project
	fmt.Printf("\n🚀 Creating project: %s\n", cfg.ProjectName)
	fmt.Printf("📦 Module: %s\n", cfg.ModuleName)
	fmt.Printf("🎯 Sample API: %s\n\n", cfg.SampleAPIName)
	
	gen := generator.NewProjectGenerator(cfg)
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}
	
	// Success message
	fmt.Printf("\n✅ Project successfully created at %s/%s\n\n", cfg.OutputDir, cfg.ProjectName)
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", cfg.ProjectName)
	fmt.Println("  make up              # Start all services")
	fmt.Println("  make migrate-up      # Run migrations")
	fmt.Println("  make run             # Start the application")
	fmt.Printf("\n🌐 Access your application at http://localhost:%d\n", cfg.AppPort)
	fmt.Printf("📚 API documentation: http://localhost:%d/api/v1/%s\n", cfg.AppPort, cfg.SampleAPINameLower+"s")
	
	return nil
}
