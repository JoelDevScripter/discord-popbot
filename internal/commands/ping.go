package commands

import (
	"project-root/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

func PingHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.Content == "!ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
		logger.Info("Ping command executed")
	}
}
