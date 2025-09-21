package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"time"
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

func LogTransaction(userID string, txType string, amount int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Pool.Exec(ctx, `
		INSERT INTO transactions (user_id, type, amount, created_at)
		VALUES ($1, $2, $3, $4)
	`, userID, txType, amount, time.Now())
	return err
}
