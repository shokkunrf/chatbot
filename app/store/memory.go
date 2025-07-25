package store

import (
	"fmt"
	"sync"
)

// ConversationHistory は特定のチャンネルのメッセージを保持する
type ConversationHistory struct {
	Messages []Message
	mu       sync.RWMutex
}

// MemoryStore はメモリ内ストレージを使用してStoreインターフェースを実装する
type MemoryStore struct {
	conversations map[string]*ConversationHistory // キー: "guildID:channelID"
	maxMessages   int
	mu            sync.RWMutex
}

// NewMemoryStore は新しいMemoryStoreインスタンスを作成する
func NewMemoryStore(maxMessages int) *MemoryStore {
	return &MemoryStore{
		conversations: make(map[string]*ConversationHistory),
		maxMessages:   maxMessages,
	}
}

// GetHistory は特定のギルドとチャンネルの会話履歴を取得する
func (s *MemoryStore) GetHistory(guildID, channelID string) []Message {
	key := fmt.Sprintf("%s:%s", guildID, channelID)

	s.mu.RLock()
	conv, exists := s.conversations[key]
	s.mu.RUnlock()

	if !exists {
		return []Message{}
	}

	conv.mu.RLock()
	defer conv.mu.RUnlock()

	// 外部からの変更を防ぐためにコピーを返す
	messages := make([]Message, len(conv.Messages))
	copy(messages, conv.Messages)
	return messages
}

// AddMessage は会話履歴に新しいメッセージを追加する
func (s *MemoryStore) AddMessage(guildID, channelID string, message Message) error {
	key := fmt.Sprintf("%s:%s", guildID, channelID)

	s.mu.Lock()
	conv, exists := s.conversations[key]
	if !exists {
		conv = &ConversationHistory{
			Messages: make([]Message, 0, s.maxMessages),
		}
		s.conversations[key] = conv
	}
	s.mu.Unlock()

	conv.mu.Lock()
	defer conv.mu.Unlock()

	// メッセージを追加
	conv.Messages = append(conv.Messages, message)

	// 必要に応じて最大メッセージ数に調整
	if len(conv.Messages) > s.maxMessages {
		// 最新のメッセージのみを保持
		conv.Messages = conv.Messages[len(conv.Messages)-s.maxMessages:]
	}

	return nil
}

// ClearHistory は特定のギルドとチャンネルの会話履歴をクリアする
func (s *MemoryStore) ClearHistory(guildID, channelID string) error {
	key := fmt.Sprintf("%s:%s", guildID, channelID)

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.conversations, key)
	return nil
}
