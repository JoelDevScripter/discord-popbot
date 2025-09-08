package main

import (
	"log"
	"os"

	"project-root/internal/commands"
	"project-root/pkg/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("Missing DISCORD_TOKEN in .env")
	}

	// Create Discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session,", err)
	}

	// Add handlers
	dg.AddHandler(commands.PingHandler)

	// Open connection
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection,", err)
	}
	defer dg.Close()

	logger.Info("Bot is running. Press CTRL+C to exit.")
	select {}
}
