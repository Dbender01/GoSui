package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dbender01/GoSui/internal/config"
)

var registeredCommands []*discordgo.ApplicationCommand
var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

func Register(s *discordgo.Session) error {
	// Add future command registrations here
	guildID := config.GetDevGuild()
	cmdList := []*discordgo.ApplicationCommand{
		registerTestCommand(),
		registerAskCommand(),
	}

	for _, cmd := range cmdList {
		created, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
		if err != nil {
			return err
		}
		registeredCommands = append(registeredCommands, created)
	}
	return nil
}

func Cleanup(s *discordgo.Session) {
	for _, cmd := range registeredCommands {
		_ = s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
	}
}

func Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}
