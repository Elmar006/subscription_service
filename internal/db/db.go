package db

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/Elmar006/subscription_service/internal/config"
	"github.com/Elmar006/subscription_service/logger"
)

func Connect(cfg *config.Config) *sql.DB {
	logger.Init()
	log := logger.L()
	db, err := sql.Open("postgres", cfg.DBConnString())
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	log.Info("Connected to PostgreSql successfully")
	return db
}
