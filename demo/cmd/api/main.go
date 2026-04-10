package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/muazwzxv/demo/cmd"
	"github.com/muazwzxv/demo/internal/config"
	"github.com/muazwzxv/demo/internal/handlers"
	"github.com/muazwzxv/demo/internal/models"
	"github.com/muazwzxv/demo/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	app := fiber.New()

	apiService := &cmd.APIService{
		DB:      db,
		Queries: models.New(),
	}

	bootupCtx := context.Background()
	handlers.SetupHandler(bootupCtx, app, apiService)

	slog.InfoContext(bootupCtx, fmt.Sprintf("Server starting on port %s", cfg.ServerPort))
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		slog.ErrorContext(bootupCtx, fmt.Sprintf("Failed to start server: %v", err))
		panic(err)
	}
}
