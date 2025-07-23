package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func NewPostgresDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://user:password@localhost:5432/subscriptions?sslmode=disable"
		log.Println("DATABASE_DSN not set, using default:", dsn)
	} else {
		log.Println("Using DATABASE_DSN from environment")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка ping к БД: %w", err)
	}

	log.Println("Успешное подключение к базе данных")
	return db, nil
}
