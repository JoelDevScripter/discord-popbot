package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Lista de fun facts cute
var funFacts = []string{
	"✨ You shine brighter than a star!",
	"💖 Your smile is contagious!",
	"🌸 Remember to sparkle today!",
	"🎀 You are pawsitively amazing!",
	"🍬 Life is sweeter with you in it!",
	"🌈 Your positivity is magical!",
	"💫 You bring light wherever you go!",
	"🦄 Believe in yourself like a unicorn!",
	"🍓 Sweetness overload detected!",
	"🐾 You leave tiny paw prints of joy!",
	"🌷 You make the world bloom!",
	"🎶 Your laugh is music to the soul!",
	"🍭 You are sugar, spice, and everything nice!",
	"💌 Sending virtual hugs your way!",
	"☀️ You brighten even the cloudiest days!",
	"💎 You sparkle in your own unique way!",
	"🌟 Stars are jealous of your shine!",
	"🐰 You are cute enough to be a mascot!",
	"🍩 Sweet like donuts, fun like sprinkles!",
	"🎉 Every day is better with your energy!",
	"🍒 Cherry-picked for awesomeness!",
	"🧸 Soft, cuddly, and full of charm!",
	"🌺 Your vibe is pure flower power!",
	"🍉 Refreshing like summer fruits!",
	"🎠 Magical moments follow you!",
	"💃 Dance like nobody's watching!",
	"🌻 Sunflowers look up just like you!",
	"🐝 Buzzing with happiness!",
	"🍦 Life is better with your sweetness!",
	"🎈 Floating on joy like a balloon!",
	"🌙 Your dreams are as big as the sky!",
	"🦋 Fluttering into everyone’s heart!",
	"🍪 You are the cookie in the jar of life!",
	"🎇 Sparkle more than fireworks tonight!",
	"💐 You are a bouquet of smiles!",
	"🐱 Cuteness level: maximum overload!",
	"🍋 Sweet, tangy, and full of zest!",
	"🌸 Petals of joy surround you!",
	"🍯 Sweet like honey, bright like gold!",
	"🦄 Believe in magic, because you are magical!",
	"🎂 Life is a cake and you’re the frosting!",
	"🌈 Your happiness is a rainbow!",
	"🐤 Chirping positivity everywhere!",
	"🍓 Strawberry sweet and berry cute!",
	"💖 Hearts follow you wherever you go!",
	"🎶 You are a melody of joy!",
	"☁️ Soft and fluffy like a cloud!",
	"🧁 Sprinkled with sweetness and charm!",
	"🌟 Twinkle twinkle little star, that's you!",
	"🎁 Wrapped in awesomeness!",
	"🐾 Little steps, big impact!",
	"🍬 Candy-coated kindness detected!",
	"💫 Your aura sparkles brighter than diamonds!",
	"🌺 Blossoming beauty all around!",
	"🦋 Your wings of joy lift spirits!",
	"🎀 Bow-tiful and fabulous!",
	"🍉 Juicy and delightful like a fruit!",
	"🎉 Confetti-worthy happiness!",
	"💌 Heart-shaped messages sent to you!",
	"🐰 Bunny-level cuteness achieved!",
	"🌞 Radiate sunshine wherever you go!",
	"🍭 Swirl of sweetness in the world!",
	"🧸 Hug dispenser activated!",
	"🌈 Your presence paints rainbows!",
	"💖 Glittering with endless charm!",
	"🎠 Carousel of smiles follows you!",
	"🍒 You are the cherry on top of life!",
	"💃 Twirl into happiness!",
	"🌻 Bloom where you are planted!",
	"🐝 Honey-level positivity engaged!",
	"🍦 Ice cream and happiness combined!",
	"🎈 Lifted by joy and laughter!",
	"🌙 Moonlight loves your glow!",
	"🦄 Magical vibes incoming!",
	"🍪 Sweet bites of joy all day!",
	"🎇 Firework of sparkle and charm!",
	"💐 Garden of smiles around you!",
	"🐱 Purr-fectly cute and lovable!",
	"🍋 Zesty happiness all around!",
	"🌸 Petals of joy everywhere!",
	"🍯 Sticky-sweet kindness!",
	"🦋 Fluttering happiness wings!",
	"🎂 Slice of joy served daily!",
	"🌈 Rainbow of positivity detected!",
	"🐤 Chirpy happiness spreading!",
	"🍓 Berries of joy sprinkled!",
	"💖 Heartbeats of sparkle!",
	"🎶 Harmony of smiles achieved!",
	"☁️ Cloud-soft sweetness!",
	"🧁 Cupcake-level adorableness!",
	"🌟 Twinkle, shine, repeat!",
	"🎁 Present of joy delivered!",
	"🐾 Footsteps of happiness follow you!",
	"🍬 Candy rainbow of delight!",
	"💫 Shining like never before!",
	"🌺 Blossom power: activated!",
	"🦋 Wings of charm deployed!",
	"🎀 Ribboned with fun!",
}

// Devuelve un fun fact aleatorio
func getRandomFunFact() string {
	rand.Seed(time.Now().UnixNano())
	return funFacts[rand.Intn(len(funFacts))]
}

// -------------------
// Función central para generar embed de perfil
// -------------------
func sendProfileEmbed(s *discordgo.Session, channelID string, user *discordgo.User, isInteraction bool, i *discordgo.Interaction) {
	// Registro automático
	if err := EnsureUserExists(user.ID); err != nil {
		msg := "Oopsie! Something went wrong accessing the profile 💖"
		if isInteraction {
			s.InteractionRespond(i, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: msg},
			})
		} else {
			s.ChannelMessageSend(channelID, msg)
		}
		return
	}

	coins, err := GetUserCoins(user.ID)
	if err != nil {
		msg := "Oopsie! Could not get Notes 💕"
		if isInteraction {
			s.InteractionRespond(i, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: msg},
			})
		} else {
			s.ChannelMessageSend(channelID, msg)
		}
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "✨ Sparkly Profile ✨",
		Description: fmt.Sprintf("Hey cutie! 💖 Here's the profile info for %s 🌸", user.Username),
		Color:       0xFFB6C1,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "💎 Notes", Value: fmt.Sprintf("%d", coins), Inline: true},
			{Name: "🆔 Discord ID", Value: user.ID, Inline: true},
			{Name: "🗓 Registered At", Value: time.Now().Format("02 Jan 2006"), Inline: false},
			{Name: "🎀 Fun Fact", Value: getRandomFunFact(), Inline: false},
		},
	}

	if isInteraction {
		s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
	} else {
		s.ChannelMessageSendEmbed(channelID, embed)
	}
}

// -------------------
// Prefijo !profile
// -------------------
func HandlePrefixProfile(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || len(m.Content) < 8 || m.Content[:8] != "!profile" {
		return
	}

	// Si el usuario menciona a alguien, usarlo
	var targetUser *discordgo.User = m.Author
	if len(m.Mentions) > 0 {
		targetUser = m.Mentions[0]
	}

	sendProfileEmbed(s, m.ChannelID, targetUser, false, nil)
}

// -------------------
// Slash command
// -------------------
func GetProfileCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "profile",
		Description: "Check your cute profile and balance 💖",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Optional: see another user's profile",
				Required:    false,
			},
		},
	}
}

func HandleSlashProfile(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var targetUser *discordgo.User

	if len(i.ApplicationCommandData().Options) > 0 {
		option := i.ApplicationCommandData().Options[0]
		if option.Type == discordgo.ApplicationCommandOptionUser && option.UserValue(s) != nil {
			targetUser = option.UserValue(s)
		}
	}

	if targetUser == nil {
		targetUser = i.Member.User
	}

	sendProfileEmbed(s, "", targetUser, true, i.Interaction)
}
