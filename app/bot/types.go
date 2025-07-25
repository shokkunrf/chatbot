package bot

import (
	"github.com/bwmarrin/discordgo"
	"chatbot/ai"
	"chatbot/store"
)

// Bot はDiscordボットのメインハンドラー
type Bot struct {
	session      *discordgo.Session
	store        store.Store
	aiClient     ai.Client
	botID        string
	systemPrompt string
}

// Config はボットの設定
type Config struct {
	Token        string
	SystemPrompt string
}
