package main

import (
	"log"
	"os"
	"os/signal"
	"project-root/internal/commands"
	"project-root/pkg/db"
	"project-root/pkg/logger"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// ----------------------------
	// Cargar variables de entorno
	// ----------------------------
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("Missing DISCORD_TOKEN in .env")
	}

	// ----------------------------
	// Conectar a la base de datos
	// ----------------------------
	err = db.Connect()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Pool.Close()

	// ----------------------------
	// Crear sesión de Discord
	// ----------------------------
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	// ----------------------------
	// Handlers de prefijo libres
	// ----------------------------
	dg.AddHandler(commands.HandlePrefixPing)
	dg.AddHandler(commands.HandlePrefixProfile)

	// ----------------------------
	// Handlers automáticos con registro obligatorio
	// ----------------------------
	commands.SetupHandlers(dg) // Aquí se aplicará la verificación para futuros comandos

	// ----------------------------
	// Abrir conexión
	// ----------------------------
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening Discord connection:", err)
	}
	defer dg.Close()

	// ----------------------------
	// Registrar todos los slash commands (limpio)
	// ----------------------------
	commands.RegisterAllCommands(dg)

	logger.Info("✅ Bot is running. Press CTRL+C to exit.")

	// ----------------------------
	// Esperar CTRL+C para cerrar
	// ----------------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("⚠️ Bot shutting down.")
}
