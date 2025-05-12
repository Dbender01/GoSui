package listeners

import (
	"log"
	"strings"
	"fmt"
	"sort"
	"github.com/bwmarrin/discordgo"
	"github.com/dbender01/GoSui/internal/ai"
	"github.com/dbender01/GoSui/internal/helpers"
)

func fetchConversationHistory(s *discordgo.Session, channelID, replyMsgID string, historyLimit int, botID string) ([]ai.Message, error) {
	stopTyping := helpers.StartTyping(s, channelID)
	defer stopTyping()

	replyMsg, err := s.ChannelMessage(channelID, replyMsgID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reply message: %w", err)
	}

	prevMessages, err := s.ChannelMessages(channelID, historyLimit - 1, replyMsgID, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch previous messages: %w", err)
	}

	allMessages := append(prevMessages, replyMsg)
	sort.Slice(allMessages, func(i, j int) bool {
		return allMessages[i].Timestamp.Before(allMessages[j].Timestamp)
	})

	conversationHistory := make([]ai.Message, 0, historyLimit)
	for i := len(allMessages) - 1; i >= 0; i-- {
		msg := allMessages[i]

		if strings.TrimSpace(msg.Content) == "" {
			continue
		}

		role := "user"
		if msg.Author.Bot {
			role = "assistant"
		}

		if msg.MessageReference == nil {
			for _, mention := range msg.Mentions {
				if mention.ID == botID {
					log.Printf("Bot mention (non-reply) detected in: %s", msg.Content)
					conversationHistory = append([]ai.Message{{
						Role:    role,
						Content: msg.Content,
					}}, conversationHistory...)
					return conversationHistory, nil
				}
			}
		}

		conversationHistory = append([]ai.Message{{
			Role:    role,
			Content: msg.Content,
		}}, conversationHistory...)
	}

	return conversationHistory, nil
}

func HandleReplyMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if m.MessageReference == nil {
		return
	}

	referencedMsg, err := s.ChannelMessage(m.ChannelID, m.MessageReference.MessageID)
	if err != nil || !referencedMsg.Author.Bot {
		return
	}

	botUser, err := s.User("@me")
	if err != nil {
		log.Printf("Error getting bot user: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Error identifying bot user.")
		return
	}
	botID := botUser.ID

	conversationHistory, err := fetchConversationHistory(s, m.ChannelID, m.MessageReference.MessageID, 10, botID)
	if err != nil {
		log.Printf("Error fetching conversation history: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Error retrieving conversation context.")
		return
	}

	conversationHistory = append(conversationHistory, ai.Message{
		Role:    "user",
		Content: m.Content,
	})
	
	log.Printf("--------Message History Retrieved: %v", len(conversationHistory))

	answer, err := ai.ContinueWithBorys(conversationHistory)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Claude had an error: "+err.Error())
		return
	}

	s.ChannelMessageSendReply(m.ChannelID, answer, m.Reference())
}
