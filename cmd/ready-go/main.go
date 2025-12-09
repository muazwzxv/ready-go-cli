package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	internalcli "github.com/muazwzxv/ready-go-cli/internal/cli"
	"github.com/muazwzxv/ready-go-cli/internal/generator"
	"github.com/urfave/cli/v2"
)

var (
	version = "1.1.0"
)

//go:embed all:templates
var embeddedTemplates embed.FS

func main() {
	// Set the embedded templates for the generator
	generator.SetEmbeddedTemplates(embeddedTemplates)

	app := &cli.App{
		Name:     "ready-go",
		Usage:    "Scaffold production-ready Go projects with clean architecture",
		Version:  version,
		Commands: internalcli.Commands(),
		Authors: []*cli.Author{
			{
				Name:  "Ready-Go CLI",
				Email: "contact@ready-go.dev",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		log.Fatal(err)
	}
}
