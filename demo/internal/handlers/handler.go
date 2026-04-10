package handlers

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/muazwzxv/demo/cmd"
	"github.com/muazwzxv/demo/internal/handlers/sample"
)

func SetupHandler(ctx context.Context, router *fiber.App, svc *cmd.APIService) {
	setupSampleHandlers(ctx, router, svc)

	// setupUserHandlers()
	// setupStoreHanlders()
}

func setupSampleHandlers(ctx context.Context, router *fiber.App, svc *cmd.APIService) {
	sampleHandler := &sample.SampleHandler{
		DB:      svc.DB,
		Queries: svc.Queries,
	}
	router.Post("/v1/demo", sampleHandler.Handle)
	slog.InfoContext(ctx, "Registered sample handlers")
}
