package helpers

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func StartTyping(s *discordgo.Session, channelID string) func() {
    stopTyping := make(chan struct{})
    
    go func() {
        err := s.ChannelTyping(channelID)
        if err != nil {
            log.Printf("Error sending typing indicator: %v", err)
        }
        
        ticker := time.NewTicker(5 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ticker.C:
                err := s.ChannelTyping(channelID)
                if err != nil {
                    log.Printf("Error refreshing typing indicator: %v", err)
                    return
                }
            case <-stopTyping:
                return
            }
        }
    }()
    
    return func() {
        close(stopTyping)
    }
}
