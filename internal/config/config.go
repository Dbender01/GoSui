package config

import "os"

func GetBotToken() string {
	return os.Getenv("DISCORD_BOT_TOKEN")
}

func GetDevGuild() string {
	return os.Getenv("TEST_GUILD_ID")
}

func GetAnthropicKey() string {
	return os.Getenv("ANTHROPIC_KEY")
}
