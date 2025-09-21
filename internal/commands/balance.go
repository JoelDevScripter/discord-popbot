package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// -------------------
// Prefijo !balance
// -------------------
func HandlePrefixBalance(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || !strings.HasPrefix(m.Content, "!balance") {
		return
	}

	// Valor por defecto: el autor
	target := m.Author

	// Si hay un mention, usamos ese usuario
	if len(m.Mentions) > 0 {
		target = m.Mentions[0]
	}

	coins, err := GetUserCoins(target.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Oopsie! Could not get their Notes 💕")
		return
	}

	embed := buildBalanceEmbed(target.Username, target.AvatarURL(""), coins)
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// -------------------
// Slash command
// -------------------
func GetBalanceCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "balance",
		Description: "Check how many Notes you (or another user) have 💎",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user you want to check (optional)",
				Required:    false,
			},
		},
	}
}

func HandleSlashBalance(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	// Valor por defecto: el usuario que ejecuta el comando
	target := i.Member.User

	// Si se pasa "user" en la opción, lo usamos
	if len(data.Options) > 0 && data.Options[0].UserValue(s) != nil {
		target = data.Options[0].UserValue(s)
	}

	coins, err := GetUserCoins(target.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Oopsie! Could not get their Notes 💕",
			},
		})
		return
	}

	embed := buildBalanceEmbed(target.Username, target.AvatarURL(""), coins)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

// -------------------
// Helper para construir embed cute
// -------------------
func buildBalanceEmbed(username, avatarURL string, coins int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "📊 Balance Report",
		Description: fmt.Sprintf("Hey here's the balance for **%s**! ✨", username),
		Color:       0xFFB6C1,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatarURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "🪙 Notes",
				Value:  fmt.Sprintf("%d", coins),
				Inline: true,
			},
			{
				Name:   "🌸 Level",
				Value:  "Coming soon…", // Placeholder
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Keep shining and collecting Notes ✨💖",
		},
	}
}
