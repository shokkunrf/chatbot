package store

import (
	"testing"
	"time"
)

// TestMemoryStore_AddAndGetMessages は基本的なメッセージの保存と取得をテスト
func TestMemoryStore_AddAndGetMessages(t *testing.T) {
	store := NewMemoryStore(20)
	guildID := "test-guild"
	channelID := "test-channel"

	// 空の履歴をテスト
	history := store.GetHistory(guildID, channelID)
	if len(history) != 0 {
		t.Errorf("Expected empty history, got %d messages", len(history))
	}

	// メッセージを追加
	msg1 := Message{
		Role:      "user",
		Content:   "Hello",
		UserID:    "user1",
		UserName:  "TestUser",
		Timestamp: time.Now(),
	}
	msg2 := Message{
		Role:      "assistant",
		Content:   "Hi there!",
		UserID:    "bot",
		UserName:  "TestBot",
		Timestamp: time.Now(),
	}

	err := store.AddMessage(guildID, channelID, msg1)
	if err != nil {
		t.Fatalf("Failed to add message 1: %v", err)
	}

	err = store.AddMessage(guildID, channelID, msg2)
	if err != nil {
		t.Fatalf("Failed to add message 2: %v", err)
	}

	// 履歴を取得
	history = store.GetHistory(guildID, channelID)
	if len(history) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(history))
	}

	// メッセージの順序を確認
	if history[0].Content != msg1.Content {
		t.Errorf("Expected first message to be '%s', got '%s'", msg1.Content, history[0].Content)
	}
	if history[1].Content != msg2.Content {
		t.Errorf("Expected second message to be '%s', got '%s'", msg2.Content, history[1].Content)
	}
}

// TestMemoryStore_MaxMessages はメッセージ数の上限制限をテスト
func TestMemoryStore_MaxMessages(t *testing.T) {
	maxMessages := 20
	store := NewMemoryStore(maxMessages)
	guildID := "test-guild"
	channelID := "test-channel"

	// 25件のメッセージを追加
	for i := 0; i < 25; i++ {
		msg := Message{
			Role:      "user",
			Content:   string(rune('A' + i)), // A, B, C, ... Y
			UserID:    "user1",
			UserName:  "TestUser",
			Timestamp: time.Now(),
		}
		err := store.AddMessage(guildID, channelID, msg)
		if err != nil {
			t.Fatalf("Failed to add message %d: %v", i, err)
		}
	}

	// 履歴を取得
	history := store.GetHistory(guildID, channelID)
	if len(history) != maxMessages {
		t.Errorf("Expected %d messages, got %d", maxMessages, len(history))
	}

	// 最新のメッセージ（FからY）を持っていることを確認
	if history[0].Content != "F" {
		t.Errorf("Expected first message to be 'F', got '%s'", history[0].Content)
	}
	if history[maxMessages-1].Content != "Y" {
		t.Errorf("Expected last message to be 'Y', got '%s'", history[maxMessages-1].Content)
	}
}

// TestMemoryStore_GuildIsolation はギルドごとに履歴が分離されていることをテスト
func TestMemoryStore_GuildIsolation(t *testing.T) {
	store := NewMemoryStore(20)
	guildA := "guild-a"
	guildB := "guild-b"
	channelA := "channel-a"
	channelB := "channel-b"

	// ギルドAにメッセージを追加
	msgA := Message{
		Role:      "user",
		Content:   "Message from Guild A",
		UserID:    "userA",
		UserName:  "UserA",
		Timestamp: time.Now(),
	}
	err := store.AddMessage(guildA, channelA, msgA)
	if err != nil {
		t.Fatalf("Failed to add message to guild A: %v", err)
	}

	// ギルドBにメッセージを追加
	msgB := Message{
		Role:      "user",
		Content:   "Message from Guild B",
		UserID:    "userB",
		UserName:  "UserB",
		Timestamp: time.Now(),
	}
	err = store.AddMessage(guildB, channelB, msgB)
	if err != nil {
		t.Fatalf("Failed to add message to guild B: %v", err)
	}

	// ギルドAの履歴を取得
	historyA := store.GetHistory(guildA, channelA)
	if len(historyA) != 1 {
		t.Errorf("Expected 1 message for guild A, got %d", len(historyA))
	}
	if historyA[0].Content != msgA.Content {
		t.Errorf("Expected guild A message content '%s', got '%s'", msgA.Content, historyA[0].Content)
	}

	// ギルドBの履歴を取得
	historyB := store.GetHistory(guildB, channelB)
	if len(historyB) != 1 {
		t.Errorf("Expected 1 message for guild B, got %d", len(historyB))
	}
	if historyB[0].Content != msgB.Content {
		t.Errorf("Expected guild B message content '%s', got '%s'", msgB.Content, historyB[0].Content)
	}

	// 相互汚染がないことを確認
	historyA2 := store.GetHistory(guildA, channelB)
	if len(historyA2) != 0 {
		t.Errorf("Expected no messages for guild A channel B, got %d", len(historyA2))
	}

	historyB2 := store.GetHistory(guildB, channelA)
	if len(historyB2) != 0 {
		t.Errorf("Expected no messages for guild B channel A, got %d", len(historyB2))
	}
}

// TestMemoryStore_ClearHistory は会話履歴のクリアをテスト
func TestMemoryStore_ClearHistory(t *testing.T) {
	store := NewMemoryStore(20)
	guildID := "test-guild"
	channelID := "test-channel"

	// いくつかメッセージを追加
	for i := 0; i < 5; i++ {
		msg := Message{
			Role:      "user",
			Content:   "Test message",
			UserID:    "user1",
			UserName:  "TestUser",
			Timestamp: time.Now(),
		}
		err := store.AddMessage(guildID, channelID, msg)
		if err != nil {
			t.Fatalf("Failed to add message: %v", err)
		}
	}

	// メッセージが存在することを確認
	history := store.GetHistory(guildID, channelID)
	if len(history) != 5 {
		t.Errorf("Expected 5 messages before clear, got %d", len(history))
	}

	// 履歴をクリア
	err := store.ClearHistory(guildID, channelID)
	if err != nil {
		t.Fatalf("Failed to clear history: %v", err)
	}

	// 履歴がクリアされたことを確認
	history = store.GetHistory(guildID, channelID)
	if len(history) != 0 {
		t.Errorf("Expected 0 messages after clear, got %d", len(history))
	}
}

// TestMemoryStore_ConcurrentAccess はスレッドセーフティをテスト
func TestMemoryStore_ConcurrentAccess(t *testing.T) {
	store := NewMemoryStore(20)
	guildID := "test-guild"
	channelID := "test-channel"

	// 並行処理を実行
	done := make(chan bool)

	// ライターgoroutine
	go func() {
		for i := 0; i < 100; i++ {
			msg := Message{
				Role:      "user",
				Content:   "Concurrent message",
				UserID:    "user1",
				UserName:  "TestUser",
				Timestamp: time.Now(),
			}
			store.AddMessage(guildID, channelID, msg)
		}
		done <- true
	}()

	// リーダーgoroutine
	go func() {
		for i := 0; i < 100; i++ {
			store.GetHistory(guildID, channelID)
		}
		done <- true
	}()

	// 両方のgoroutineの完了を待つ
	<-done
	<-done

	// パニックなしでここまで来れば、テスト成功
	// メッセージがあることを確認
	history := store.GetHistory(guildID, channelID)
	if len(history) == 0 {
		t.Error("Expected some messages after concurrent access")
	}
}
