package commands

import (
	"project-root/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

// Estructura para un comando slash
type SlashCommand struct {
	Definition *discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// Lista global de comandos
var slashCommands = []SlashCommand{
	{GetPingCommand(), HandleSlashPing},
	{GetProfileCommand(), HandleSlashProfile},
}

// SetupHandlers conecta los handlers automáticamente
func SetupHandlers(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}
		for _, cmd := range slashCommands {
			if i.ApplicationCommandData().Name == cmd.Definition.Name {
				cmd.Handler(s, i)
			}
		}
	})
}

// RegisterAllCommands borra y vuelve a registrar todos los slash commands
func RegisterAllCommands(s *discordgo.Session) {
	// 1. Borrar comandos viejos
	existing, _ := s.ApplicationCommands(s.State.User.ID, "")
	for _, cmd := range existing {
		_ = s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
	}

	// 2. Registrar nuevos
	for _, cmd := range slashCommands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd.Definition)
		if err != nil {
			logger.Error("Error registering command:", err)
		}
	}
}
