package config

import "os"

func GetBotToken() string {
	return os.Getenv("DISCORD_BOT_TOKEN")
}