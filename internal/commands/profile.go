package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Prefijo !profile
func HandlePrefixProfile(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content != "!profile" {
		return
	}

	err := EnsureUserExists(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error accessing your profile.")
		return
	}

	coins, err := GetUserCoins(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error getting your coins.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, you have %d Notes.", m.Author.Username, coins))
}

// Definición para el slash command
func GetProfileCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "profile",
		Description: "Check your profile and balance.",
	}
}

// Handler del slash
func HandleSlashProfile(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "profile" {
		return
	}

	user := i.Member.User

	err := EnsureUserExists(user.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error accessing your profile.",
			},
		})
		return
	}

	coins, err := GetUserCoins(user.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error getting your coins.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s, you have %d Notes.", user.Username, coins),
		},
	})
}
