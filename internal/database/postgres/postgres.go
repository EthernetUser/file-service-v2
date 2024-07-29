package postgres

import (
	"database/sql"
	"file-service/m/internal/config"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func New(cfg config.DatabaseConfig) *sql.DB {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err := sql.Open("postgres", connString)

	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	return db
}
