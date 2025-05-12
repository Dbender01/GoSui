package listeners

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dbender01/GoSui/internal/ai"
	"github.com/dbender01/GoSui/internal/helpers"
)

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore bot messages
	if m.Author.Bot {
		return
	}

	// Check if this is a reply - if so, ignore it as HandleReplyMessages will handle it
	if m.MessageReference != nil {
    	return
	}

	// Check if the bot is mentioned
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

	// Remove the mention prefix so we get just the question
	content := strings.TrimSpace(strings.Replace(m.Content, "<@"+botID+">", "", 1))

	// Now respond using AI

	response, err := ai.AskBorys(content)

	if err != nil {
		response = "Claude had an error: " + err.Error()
	}

	s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
}
