package commands

import (
	"context"
	"fmt"
	"project-root/pkg/db"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Revisar si el usuario está registrado
func IsUserRegistered(discordID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	err := db.Pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE discord_id=$1)", discordID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking user: %w", err)
	}

	return exists, nil
}

// Middleware para comandos de prefijo
func RequireRegistrationPrefix(s *discordgo.Session, m *discordgo.MessageCreate, commandFunc func()) {
	if m.Content == "!ping" || m.Content == "!profile" {
		commandFunc()
		return
	}

	registered, err := IsUserRegistered(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Oops! Something went wrong checking your registration 💖")
		return
	}

	if !registered {
		s.ChannelMessageSend(m.ChannelID, "Hey cutie! 💕 You need to use `!profile` first to get registered ✨")
		return
	}

	commandFunc()
}

// Middleware para slash commands
func RequireRegistrationSlash(s *discordgo.Session, i *discordgo.InteractionCreate, commandFunc func()) {
	name := i.ApplicationCommandData().Name
	if name == "ping" || name == "profile" {
		commandFunc()
		return
	}

	registered, err := IsUserRegistered(i.Member.User.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Oops! Something went wrong checking your registration 💖",
			},
		})
		return
	}

	if !registered {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey cutie! 💕 You need to use `/profile` first to get registered ✨",
			},
		})
		return
	}

	commandFunc()
}
