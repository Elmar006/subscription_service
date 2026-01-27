package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/Elmar006/subscription_service/internal/config"
)

func Connect(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.DBConnString())
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	log.Println("Connected to PostgreSql successfully")
	return db
}
