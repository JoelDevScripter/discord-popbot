package commands

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"project-root/pkg/db"
	"project-root/pkg/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"
)

const (
	MinDailyNotes = 50
	MaxDailyNotes = 150
)

// HandlePrefixDaily for !daily
func HandlePrefixDaily(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	if m.Content != "!daily" {
		return
	}

	RequireRegistrationPrefix(s, m, func() {
		handleDaily(s, m.Author, m.ChannelID, false, nil)
	})
}

// GetDailyCommand defines the /daily slash command
func GetDailyCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "daily",
		Description: "Claim your daily Notes ✨",
	}
}

// HandleSlashDaily for /daily
func HandleSlashDaily(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name != "daily" {
		return
	}

	RequireRegistrationSlash(s, i, func() {
		// i.Member should exist for guild commands
		if i.Member == nil || i.Member.User == nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Could not determine the user (try this command in a server).",
				},
			})
			return
		}
		handleDaily(s, i.Member.User, "", true, i.Interaction)
	})
}

// handleDaily central logic (works for both prefix and slash)
func handleDaily(s *discordgo.Session, user *discordgo.User, channelID string, isInteraction bool, inter *discordgo.Interaction) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start transaction (pgx.TxOptions{} for defaults)
	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		logger.Error("BeginTx error:", err)
		reply(s, channelID, isInteraction, inter, "Oops! Could not start the transaction 💔")
		return
	}
	// Ensure rollback if something fails; commit explicitly later.
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// We allow last_daily to be NULL -> scan into *time.Time
	var lastClaim *time.Time
	err = tx.QueryRow(ctx, "SELECT last_daily FROM users WHERE discord_id=$1", user.ID).Scan(&lastClaim)
	if err != nil {
		// If user row doesn't exist -> ask to register
		if err == pgx.ErrNoRows {
			reply(s, channelID, isInteraction, inter, "You need to register first using `/profile` 💖")
			return
		}
		logger.Error("QueryRow last_daily error:", err)
		reply(s, channelID, isInteraction, inter, "Oops! Something went wrong accessing your Notes 💖")
		return
	}

	now := time.Now().UTC()

	// If lastClaim is not nil, check cooldown
	if lastClaim != nil {
		nextAvailable := lastClaim.Add(24 * time.Hour)
		if now.Before(nextAvailable) {
			remaining := nextAvailable.Sub(now)
			h := int(remaining.Hours())
			m := int(remaining.Minutes()) % 60
			reply(s, channelID, isInteraction, inter,
				fmt.Sprintf("⏳ You already claimed your daily! Come back in %dh %dm 💕", h, m))
			return
		}
	}

	// Generate reward using a local RNG (no global Seed call)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	reward := rnd.Intn(MaxDailyNotes-MinDailyNotes+1) + MinDailyNotes

	// Update user coins and last_daily, returning new coins
	var newCoins int
	err = tx.QueryRow(ctx,
		"UPDATE users SET coins = coins + $1, last_daily = $2 WHERE discord_id = $3 RETURNING coins",
		reward, now, user.ID).Scan(&newCoins)
	if err != nil {
		logger.Error("Update users error:", err)
		reply(s, channelID, isInteraction, inter, "Oops! Something went wrong updating your balance 💖")
		return
	}

	// Insert transaction (assumes transactions table stores user_id referencing users.id)
	_, err = tx.Exec(ctx,
		"INSERT INTO transactions (user_id, type, amount, created_at) VALUES ((SELECT id FROM users WHERE discord_id=$1), $2, $3, $4)",
		user.ID, "daily", reward, now)
	if err != nil {
		logger.Error("Insert transaction error:", err)
		reply(s, channelID, isInteraction, inter, "Oops! Could not save your transaction 💔")
		return
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		logger.Error("Commit error:", err)
		reply(s, channelID, isInteraction, inter, "Oops! Failed to finalize your daily 💖")
		return
	}

	// Build embed response
	embed := &discordgo.MessageEmbed{
		Title:       "💎 Daily Notes Claimed!",
		Description: fmt.Sprintf("Hey %s! You received **%d Notes** today! 🌸", user.Username, reward),
		Color:       0xFFB6C1,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL(""),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Come back every 24h for more Notes!",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "🪙 New Balance",
				Value:  fmt.Sprintf("%d Notes", newCoins),
				Inline: true,
			},
		},
	}

	// Send embed
	if isInteraction {
		_ = s.InteractionRespond(inter, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
		})
	} else {
		_, err := s.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			logger.Error("ChannelMessageSendEmbed error:", err)
		}
	}
}

// reply helper sends the same message either as interaction response or channel message
func reply(s *discordgo.Session, channelID string, isInteraction bool, inter *discordgo.Interaction, content string) {
	if isInteraction {
		_ = s.InteractionRespond(inter, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Content: content},
		})
	} else {
		_, err := s.ChannelMessageSend(channelID, content)
		if err != nil {
			logger.Error("ChannelMessageSend error:", err)
		}
	}
}
