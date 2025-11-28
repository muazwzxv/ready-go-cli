package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/muazwzxv/ready-go-cli/internal/config"
)

// PromptForConfig prompts the user for project configuration
func PromptForConfig(cfg *config.ProjectConfig) error {
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Println("\n📋 Let's configure your project!\n")
	
	// Module name
	moduleName, err := prompt(reader, fmt.Sprintf("Module name (%s)", cfg.ModuleName), cfg.ModuleName)
	if err != nil {
		return err
	}
	if moduleName != "" {
		cfg.ModuleName = moduleName
	}
	
	// Description
	description, err := prompt(reader, fmt.Sprintf("Project description (%s)", cfg.Description), cfg.Description)
	if err != nil {
		return err
	}
	if description != "" {
		cfg.Description = description
	}
	
	// Author
	author, err := prompt(reader, "Author name (optional)", "")
	if err != nil {
		return err
	}
	cfg.Author = author
	
	// Sample API name
	sampleAPI, err := prompt(reader, fmt.Sprintf("Sample API entity name (%s)", cfg.SampleAPIName), cfg.SampleAPIName)
	if err != nil {
		return err
	}
	if sampleAPI != "" {
		cfg.SampleAPIName = sampleAPI
	}
	
	// Redis
	withRedis, err := promptBool(reader, "Include Redis?", cfg.WithRedis)
	if err != nil {
		return err
	}
	cfg.WithRedis = withRedis
	
	// Kafka
	withKafka, err := promptBool(reader, "Include Kafka?", cfg.WithKafka)
	if err != nil {
		return err
	}
	cfg.WithKafka = withKafka
	
	// Skip Git
	skipGit, err := promptBool(reader, "Skip git initialization?", cfg.SkipGit)
	if err != nil {
		return err
	}
	cfg.SkipGit = skipGit
	
	fmt.Println()
	return nil
}

// prompt prompts the user for a string value
func prompt(reader *bufio.Reader, message, defaultValue string) (string, error) {
	if defaultValue != "" {
		fmt.Printf("? %s: ", message)
	} else {
		fmt.Printf("? %s: ", message)
	}
	
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}
	
	return input, nil
}

// promptBool prompts the user for a boolean value
func promptBool(reader *bufio.Reader, message string, defaultValue bool) (bool, error) {
	defaultStr := "N"
	if defaultValue {
		defaultStr = "Y"
	}
	
	fmt.Printf("? %s (y/N) [%s]: ", message, defaultStr)
	
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	
	input = strings.TrimSpace(strings.ToLower(input))
	
	if input == "" {
		return defaultValue, nil
	}
	
	return input == "y" || input == "yes", nil
}
