package store

import "time"

// Message は会話内の単一のメッセージを表す
type Message struct {
	Role      string    `json:"role"` // "user" または "assistant"
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Timestamp time.Time `json:"timestamp"`
}

// Store インターフェースは会話履歴を管理するメソッドを定義する
type Store interface {
	GetHistory(guildID, channelID string) []Message
	AddMessage(guildID, channelID string, message Message) error
	ClearHistory(guildID, channelID string) error
}
