package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dbender01/GoSui/internal/commands"
	"github.com/dbender01/GoSui/internal/config"
)

func main() {
	token := config.GetBotToken()

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Failed to create bot session: %v", err)
	}

	dg.AddHandler(commands.Handler)

	err = dg.Open()
	if err != nil {
		log.Fatalf("Failed to open Discord session: %v", err)
	}
	defer dg.Close()

	if err := commands.Register(dg); err != nil {
		log.Fatalf("Failed to register commands: %v", err)
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	commands.Cleanup(dg)
}
