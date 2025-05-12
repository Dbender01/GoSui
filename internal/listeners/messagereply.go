package listeners

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dbender01/GoSui/internal/ai"
	"github.com/dbender01/GoSui/internal/helpers"
)

// fetchConversationHistory retrieves the conversation history for a specific thread
// starting from the last @ mention of the bot or most recent 10 messages
func fetchConversationHistory(s *discordgo.Session, channelID string, replyMsgID string, limit int, botID string) ([]ai.Message, error) {
	
	stopTyping := helpers.StartTyping(s, channelID)
	defer stopTyping() 
	
	// Fetch messages in the channel
	messages, err := s.ChannelMessages(channelID, limit, "", "", replyMsgID)
	if err != nil {
		return nil, err
	}

	// Prepare conversation history
	conversationHistory := make([]ai.Message, 0, len(messages))
	
	// Flag to track whether we've found the initial @ mention
	foundInitialMention := false
	
	// Iterate through messages in reverse to build conversation context
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		
		// Determine the role based on whether the message is from a bot or a user
		role := "user"
		if msg.Author.Bot {
			role = "assistant"
		}

		// Skip empty messages
		if strings.TrimSpace(msg.Content) == "" {
			log.Println("Skipping empty message")
			continue
		}

		// Check if this message contains a mention of our bot
		isBotMention := false
		for _, mention := range msg.Mentions {
			if mention.ID == botID {
				isBotMention = true
				break
			}
		}

		// If this is a bot mention and we haven't found the initial mention yet,
		// mark this as the start of the conversation
		if isBotMention && !foundInitialMention && role == "user" {
			foundInitialMention = true
			log.Printf("Found initial bot mention: %s", msg.Content)
		}

		// Skip messages before the initial bot mention
		if !foundInitialMention {
			continue
		}

		// Add message to conversation history
		conversationHistory = append(conversationHistory, ai.Message{
			Role:    role,
			Content: msg.Content,
		})

		// Stop if we've reached the reference message (original bot message)
		if msg.ID == replyMsgID {
			break
		}
	}

	// If we didn't find a mention, use all messages we have
	if !foundInitialMention {
		log.Println("No bot mention found, using all available messages")
		conversationHistory = []ai.Message{}
		
		for i := len(messages) - 1; i >= 0; i-- {
			msg := messages[i]
			
			role := "user"
			if msg.Author.Bot {
				role = "assistant"
			}
			
			if strings.TrimSpace(msg.Content) != "" {
				conversationHistory = append(conversationHistory, ai.Message{
					Role:    role,
					Content: msg.Content,
				})
			}
			
			if msg.ID == replyMsgID {
				break
			}
		}
	}

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