package helpers

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// startTyping starts the typing indicator and returns a function to stop it
func StartTyping(s *discordgo.Session, channelID string) func() {
    // Start a goroutine to keep the typing indicator active
    stopTyping := make(chan struct{})
    
    go func() {
        // Send initial typing indicator
        err := s.ChannelTyping(channelID)
        if err != nil {
            log.Printf("Error sending typing indicator: %v", err)
        }
        
        ticker := time.NewTicker(5 * time.Second) // Discord's typing indicator lasts ~10 seconds, refresh every 5
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                // Refresh typing indicator
                err := s.ChannelTyping(channelID)
                if err != nil {
                    log.Printf("Error refreshing typing indicator: %v", err)
                    return
                }
            case <-stopTyping:
                // Stop typing indicator loop
                return
            }
        }
    }()
    
    // Return a function that will stop the typing indicator
    return func() {
        close(stopTyping)
    }
}