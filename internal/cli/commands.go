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
		AddCommand(),
	}
}

// AddCommand creates the 'add' command with subcommands
func AddCommand() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "Add components to an existing project",
		Subcommands: []*cli.Command{
			EntitySubcommand(),
		},
	}
}

// EntitySubcommand creates the 'entity' subcommand
func EntitySubcommand() *cli.Command {
	return &cli.Command{
		Name:      "entity",
		Usage:     "Add a new entity with migration to an existing project",
		ArgsUsage: "<entity-name>",
		Action:    addEntityAction,
	}
}

// addEntityAction handles the 'add entity' command execution
func addEntityAction(c *cli.Context) error {
	entityName := c.Args().First()

	if entityName == "" {
		return fmt.Errorf("entity name is required\nUsage: ready-go add entity <EntityName>")
	}

	cfg := config.NewEntityConfig(entityName)
	cfg.Process()

	if err := cfg.Validate(); err != nil {
		return err
	}

	fmt.Printf("\n🔍 Detected project at: %s\n", cfg.ProjectPath)
	fmt.Printf("🚀 Adding entity: %s\n\n", cfg.EntityName)

	gen := generator.NewEntityGenerator(cfg)
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("failed to generate entity: %w", err)
	}

	fmt.Printf("\n✅ Entity '%s' added successfully!\n\n", cfg.EntityName)
	fmt.Println("Next steps:")
	fmt.Println("  1. Edit the migration file to customize your table schema")
	fmt.Println("  2. Add SQLC queries to database/queries/")
	fmt.Println("  3. Run: make sqlc-generate")
	fmt.Println("  4. Run: make migrate-up")

	return nil
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
				Name:  "port",
				Usage: "Server port",
				Value: "8080",
			},
			&cli.StringFlag{
				Name:  "db-port",
				Usage: "MySQL port",
				Value: "3306",
			},
			&cli.StringFlag{
				Name:  "redis-port",
				Usage: "Redis port",
				Value: "6379",
			},
			&cli.StringFlag{
				Name:  "kafka-port",
				Usage: "Kafka port",
				Value: "9092",
			},
			&cli.StringFlag{
				Name:  "sample-name",
				Usage: "Sample entity name",
				Value: "User",
			},
		},
		Action: newProjectAction,
	}
}

// newProjectAction handles the 'new' command execution
func newProjectAction(c *cli.Context) error {
	projectName := c.Args().First()

	if projectName == "" {
		return fmt.Errorf("project name is required\nUsage: ready-go new [flags] <project-name>")
	}

	cfg := config.NewProjectConfig(projectName)

	// Apply flags
	if c.IsSet("module") {
		cfg.ModuleName = c.String("module")
	}
	if c.IsSet("port") {
		cfg.ServerPort = c.String("port")
	}
	if c.IsSet("db-port") {
		cfg.DBPort = c.String("db-port")
	}
	if c.IsSet("redis-port") {
		cfg.RedisPort = c.String("redis-port")
	}
	if c.IsSet("kafka-port") {
		cfg.KafkaPort = c.String("kafka-port")
	}
	if c.IsSet("sample-name") {
		cfg.SampleAPIName = c.String("sample-name")
	}

	cfg.Process()

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	fmt.Printf("\n🚀 Creating project: %s\n", cfg.ProjectName)
	fmt.Printf("📦 Module: %s\n", cfg.ModuleName)
	fmt.Printf("🎯 Sample API: %s\n\n", cfg.SampleAPIName)

	gen := generator.NewProjectGenerator(cfg)
	if err := gen.Generate(); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	fmt.Printf("\n✅ Project successfully created at ./%s\n\n", cfg.ProjectName)
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", cfg.ProjectName)
	fmt.Println("  make docker-up      # Start all services")
	fmt.Println("  make migrate-up     # Run migrations")
	fmt.Println("  make sqlc-generate  # Generate SQLC models")
	fmt.Println("  make run-api        # Start the application")
	fmt.Printf("\n🌐 Access your application at http://localhost:%s\n", cfg.ServerPort)

	return nil
}
