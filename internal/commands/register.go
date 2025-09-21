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
	{GetBalanceCommand(), HandleSlashBalance},
	{GetDailyCommand(), HandleSlashDaily},
	// Nuevos comandos: solo agregas aquí
}

// SetupHandlers conecta los handlers automáticamente y aplica registro obligatorio
func SetupHandlers(s *discordgo.Session) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		name := i.ApplicationCommandData().Name
		// Ping y Profile no requieren registro
		if name != "ping" && name != "profile" {
			RequireRegistrationSlash(s, i, func() {
				for _, cmd := range slashCommands {
					if i.ApplicationCommandData().Name == cmd.Definition.Name {
						cmd.Handler(s, i)
					}
				}
			})
			return
		}

		for _, cmd := range slashCommands {
			if name == cmd.Definition.Name {
				cmd.Handler(s, i)
			}
		}
	})
}

// RegisterAllCommands registra los slash commands de manera inteligente
func RegisterAllCommands(s *discordgo.Session) {
	existing, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		logger.Error("Error fetching existing commands:", err)
		return
	}

	for _, cmd := range slashCommands {
		found := false
		for _, e := range existing {
			if e.Name == cmd.Definition.Name {
				found = true

				// Opcional: actualizar comando si descripción cambió
				if e.Description != cmd.Definition.Description {
					_, err := s.ApplicationCommandEdit(s.State.User.ID, "", e.ID, cmd.Definition)
					if err != nil {
						logger.Error("Error updating command:", err)
					} else {
						logger.Info("Updated command:", cmd.Definition.Name)
					}
				}

				break
			}
		}

		if !found {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd.Definition)
			if err != nil {
				logger.Error("Error registering command:", err)
			} else {
				logger.Info("Registered new command:", cmd.Definition.Name)
			}
		}
	}
}
