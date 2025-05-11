package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dbender01/GoSui/internal/mindpalace"
)

func registerTestCommand() *discordgo.ApplicationCommand {
	cmd := &discordgo.ApplicationCommand{
		Name:        "test",
		Description: "Replies with a test message",
	}

	commandHandlers[cmd.Name] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		response, err := mindpalace.RandomResponseFromFile("borys.json")
		if err != nil {
			response = "I couldn't find that in my mind palace"
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
	}

	return cmd
}
