package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	//Creates discord session
	Token := os.Getenv("DISCORD_BOT_TOKEN")
	if Token == "" {
		fmt.Println("DISCORD_BOT_TOKEN not set")
		return
	}

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating discord session: ", err)
		return
	}
	// Register the message handler
	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore messages from the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Simple response
		if m.Content == "!hello" {
			s.ChannelMessageSend(m.ChannelID, "Hello, world!")
		}
	})

	// Open the websocket connection to Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return
	}
	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	// Wait for a termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	fmt.Println("Shutting down bot...")
}
