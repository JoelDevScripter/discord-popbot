package commands

import (
	"fmt"
	"project-root/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

const PingCommandName = "ping"

// Prefijo !ping
func HandlePrefixPing(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.Content == "!ping" {
		latency := s.HeartbeatLatency().Milliseconds()
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("🏓 Pong! Latency: %dms", latency))
		if err != nil {
			logger.Error("Error sending Pong:", err)
		} else {
			logger.Info("!ping executed by", m.Author.Username)
		}
	}
}

// Definición para el slash command
func GetPingCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        PingCommandName,
		Description: "Replies with Pong! and shows latency",
	}
}

// Handler del slash
func HandleSlashPing(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != PingCommandName {
		return
	}

	latency := s.HeartbeatLatency().Milliseconds()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🏓 Pong! Latency: %dms", latency),
		},
	})
	if err != nil {
		logger.Error("Error responding to /ping:", err)
	} else {
		logger.Info("/ping executed by", i.Member.User.Username)
	}
}
