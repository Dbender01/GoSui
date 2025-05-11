package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/dbender01/GoSui/internal/ai"
)

func registerAskCommand() *discordgo.ApplicationCommand {

	commandHandlers["ask"] = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		question := i.ApplicationCommandData().Options[0].StringValue()

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		log.Println(question)

		answer, err := ai.AskAnthropic(question)
		if err != nil {
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Content: "Claude had an error: " + err.Error(),
			})
			return
		}

		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: answer,
		})
	}

	cmd := &discordgo.ApplicationCommand{
		Name:        "ask",
		Description: "Ask Claude (Anthropic AI) a question.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "question",
				Description: "The question to ask Claude.",
				Required:    true,
			},
		},
	}

	return cmd
}
