package cmd

import (
	"database/sql"

	"github.com/muazwzxv/demo/internal/models"
)

type APIService struct {
	DB      *sql.DB
	Queries *models.Queries
}
