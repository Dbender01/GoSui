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

// fetchConversationHistory retrieves the conversation history for a specific thread
// starting from the last @ mention of the bot or most recent 10 messages
func fetchConversationHistory(s *discordgo.Session, channelID, replyMsgID string, historyLimit int, botID string) ([]ai.Message, error) {
	stopTyping := helpers.StartTyping(s, channelID)
	defer stopTyping()

	// Step 1: Get the reply message
	replyMsg, err := s.ChannelMessage(channelID, replyMsgID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reply message: %w", err)
	}

	// Step 2: Get up to 9 messages before it
	prevMessages, err := s.ChannelMessages(channelID, historyLimit - 1, replyMsgID, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch previous messages: %w", err)
	}

	// Step 3: Combine and reverse (so oldest to newest)
	allMessages := append(prevMessages, replyMsg)
	sort.Slice(allMessages, func(i, j int) bool {
		return allMessages[i].Timestamp.Before(allMessages[j].Timestamp)
	})

	// Step 4: Build conversation
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

		// Check for @mention, and make sure this is not just a reply in a chain
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

	// Return all if no mention was found
	return conversationHistory, nil
}

// HandleReplyMessages handles replies to the bot's messages
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

	// Get the bot's own user ID
	botUser, err := s.User("@me")
	if err != nil {
		log.Printf("Error getting bot user: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Error identifying bot user.")
		return
	}
	botID := botUser.ID

	// Fetch conversation history (limit to last 10 messages to avoid excessive context)
	conversationHistory, err := fetchConversationHistory(s, m.ChannelID, m.MessageReference.MessageID, 10, botID)
	if err != nil {
		log.Printf("Error fetching conversation history: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Error retrieving conversation context.")
		return
	}

	// Add the current reply to the conversation history
	conversationHistory = append(conversationHistory, ai.Message{
		Role:    "user",
		Content: m.Content,
	})
	
	log.Printf("--------Message History Retrieved: %v", len(conversationHistory))

	// Call the AI with the entire conversation context
	answer, err := ai.ContinueWithBorys(conversationHistory)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Claude had an error: "+err.Error())
		return
	}

	// Send response
	s.ChannelMessageSendReply(m.ChannelID, answer, m.Reference())
}

// HandleMentions handles direct mentions of the bot
// func HandleMentions(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	// Ignore bot messages
// 	if m.Author.Bot {
// 		return
// 	}
	
// 	// Get bot user
// 	botUser, err := s.User("@me")
// 	if err != nil {
// 		log.Printf("Error getting bot user: %v", err)
// 		return
// 	}
	
// 	// Check if the bot was mentioned
// 	botMentioned := false
// 	for _, mention := range m.Mentions {
// 		if mention.ID == botUser.ID {
// 			botMentioned = true
// 			break
// 		}
// 	}
	
// 	if !botMentioned {
// 		return
// 	}
	
// 	// This is a new conversation started with a mention
// 	conversationHistory := []ai.Message{
// 		{
// 			Role:    "user",
// 			Content: m.Content,
// 		},
// 	}

// 	// Call the AI with this message
// 	answer, err := ai.ContinueWithBorys(conversationHistory)
// 	if err != nil {
// 		s.ChannelMessageSend(m.ChannelID, "Claude had an error: "+err.Error())
// 		return
// 	}
	
// 	// Send response
// 	s.ChannelMessageSendReply(m.ChannelID, answer, m.Reference())
// }