package sample

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/muazwzxv/demo/internal/models"
)

type SampleHandler struct {
	DB      *sql.DB
	Queries *models.Queries
}

func (h *SampleHandler) Handle(ctx *fiber.Ctx) error {
	return nil
}
