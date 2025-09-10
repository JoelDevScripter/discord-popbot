package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() error {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return fmt.Errorf("DATABASE_URL not set in .env")
	}

	var err error
	Pool, err = pgxpool.New(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}

	fmt.Println("✅ Connected to database")
	return nil
}
