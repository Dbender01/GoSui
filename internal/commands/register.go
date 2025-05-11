package commands

import (
	"github.com/bwmarrin/discordgo"
)

var registeredCommands []*discordgo.ApplicationCommand
var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

func Register(s *discordgo.Session) error {
	// Add future command registrations here
	cmdList := []*discordgo.ApplicationCommand{
		registerTestCommand(),
	}

	for _, cmd := range cmdList {
		created, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
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
