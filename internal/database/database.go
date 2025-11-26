package database

import (
	"log/slog"

	"github.com/NKV510/question-answer-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDataBase(cfg *config.Config) error {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return err
	}

	slog.Info("connection to the database has been established")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
