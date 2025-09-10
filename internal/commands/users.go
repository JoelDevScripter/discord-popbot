package commands

import (
	"context"
	"fmt"
	"project-root/pkg/db"
	"time"
)

// EnsureUserExists revisa si el usuario está en la DB y lo crea si no.
func EnsureUserExists(discordID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	err := db.Pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE discord_id=$1)",
		discordID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking user: %w", err)
	}

	if !exists {
		// Inserta usuario con 100 monedas iniciales
		_, err := db.Pool.Exec(ctx,
			"INSERT INTO users(discord_id, coins, created_at) VALUES($1, $2, $3)",
			discordID, 100, time.Now())
		if err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}
	}

	return nil
}

// GetUserCoins devuelve las monedas del usuario
func GetUserCoins(discordID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var coins int
	err := db.Pool.QueryRow(ctx, "SELECT coins FROM users WHERE discord_id=$1", discordID).Scan(&coins)
	if err != nil {
		return 0, fmt.Errorf("error getting user coins: %w", err)
	}

	return coins, nil
}
