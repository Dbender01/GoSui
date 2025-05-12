package listeners

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dbender01/GoSui/internal/ai"
	"github.com/dbender01/GoSui/internal/helpers"
)

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.MessageReference != nil {
    	return
	}

	botID := s.State.User.ID
	mentioned := false
	for _, user := range m.Mentions {
		if user.ID == botID {
			mentioned = true
			break
		}
	}
	if !mentioned {
		return
	}
	
	stopTyping := helpers.StartTyping(s, m.ChannelID)
	defer stopTyping() 

	content := strings.TrimSpace(strings.Replace(m.Content, "<@"+botID+">", "", 1))

	response, err := ai.AskBorys(content)

	if err != nil {
		response = "Claude had an error: " + err.Error()
	}

	s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
}
