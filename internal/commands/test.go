package commands

import (
	"github.com/bwmarrin/discordgo"
)

func registerTestCommand() *discordgo.ApplicationCommand {
	cmd := &discordgo.ApplicationCommand{
		Name:        "test",
		Description: "Replies with a test message",
	}

	commandHandlers[cmd.Name] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "✅ This is a test response!",
			},
		})
	}

	return cmd
}
