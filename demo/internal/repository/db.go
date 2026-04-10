package repository

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/muazwzxv/demo/internal/config"
)

func NewDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
