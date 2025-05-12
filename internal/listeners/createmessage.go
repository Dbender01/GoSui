package listeners

import (
	"github.com/bwmarrin/discordgo"

	"github.com/dbender01/GoSui/internal/ai"
)

func HandleReplyMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot messages
	if m.Author.Bot {
		return
	}

	// Only respond to replies
	if m.MessageReference == nil {
		return
	}

	// Check if it's replying to the bot
	referencedMsg, err := s.ChannelMessage(m.ChannelID, m.MessageReference.MessageID)
	if err != nil || !referencedMsg.Author.Bot {
		return
	}

	// Call the AI with the content of the reply
	question := m.Content
	answer, err := ai.AskBorys(question)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Claude had an error: "+err.Error())
		return
	}

	// Send response
	s.ChannelMessageSendReply(m.ChannelID, answer, m.Reference())
}