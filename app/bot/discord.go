package bot

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"chatbot/ai"
	"chatbot/store"
)

// NewBot は新しいボットインスタンスを作成する
func NewBot(config Config, store store.Store, aiClient ai.Client) (*Bot, error) {
	if config.Token == "" {
		return nil, errors.New("Discordトークンが必要です")
	}

	// Discord セッションを作成
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, fmt.Errorf("Discordセッションの作成に失敗しました: %w", err)
	}

	bot := &Bot{
		session:      session,
		store:        store,
		aiClient:     aiClient,
		systemPrompt: config.SystemPrompt,
	}

	// イベントハンドラーを登録
	session.AddHandler(bot.messageCreate)
	session.AddHandler(bot.ready)
	session.AddHandler(bot.interactionCreate)

	// インテントを設定（メッセージ関連のみ）
	session.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsMessageContent

	return bot, nil
}

// Start はボットを起動する
func (b *Bot) Start() error {
	err := b.session.Open()
	if err != nil {
		return fmt.Errorf("Discord接続の開始に失敗しました: %w", err)
	}

	// スラッシュコマンドを登録
	_, err = b.session.ApplicationCommandCreate(b.session.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "reset",
		Description: "このチャンネルの会話履歴をリセットします",
	})
	if err != nil {
		log.Printf("スラッシュコマンドの登録に失敗しました: %v", err)
	}

	return nil
}

// Stop はボットを停止する
func (b *Bot) Stop() error {
	return b.session.Close()
}

// ready はボットが準備完了した時のハンドラー
func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	b.botID = event.User.ID
	log.Printf("ボットが起動しました: %s#%s", event.User.Username, event.User.Discriminator)
}

// messageCreate は新しいメッセージが作成された時のハンドラー
func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 自分のメッセージは無視
	if m.Author.ID == s.State.User.ID {
		return
	}

	// 応答すべきかチェック
	if !b.shouldRespond(m.Message) {
		return
	}

	// タイピング表示
	s.ChannelTyping(m.ChannelID)

	// メッセージを処理
	response, err := b.processMessage(m.Message)
	if err != nil {
		log.Printf("メッセージ処理エラー: %v", err)
		response = "申し訳ありません。現在応答できません。しばらくしてからお試しください。"
	}

	// リプライとして応答を送信
	_, err = s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
	if err != nil {
		log.Printf("メッセージ送信エラー: %v", err)
	}
}

// interactionCreate はインタラクション（スラッシュコマンド）が作成された時のハンドラー
func (b *Bot) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "reset" {
		// リセットコマンドを処理
		err := b.handleResetCommand(i.GuildID, i.ChannelID)

		var content string
		if err != nil {
			log.Printf("リセットコマンドエラー: %v", err)
			content = "エラーが発生しました。管理者にお問い合わせください。"
		} else {
			content = "このチャンネルの会話履歴をリセットしました。"
		}

		// インタラクションに応答
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
		if err != nil {
			log.Printf("インタラクション応答エラー: %v", err)
		}
	}
}

// shouldRespond はメッセージに応答すべきかどうかを判定する
func (b *Bot) shouldRespond(m *discordgo.Message) bool {
	// ボットからのメッセージは無視
	if m.Author.Bot {
		return false
	}

	// DMは無視
	if m.GuildID == "" {
		return false
	}

	// ボットへのメンションをチェック
	for _, mention := range m.Mentions {
		if mention.ID == b.botID {
			return true
		}
	}

	// ボットへの返信をチェック
	if m.ReferencedMessage != nil && m.ReferencedMessage.Author.ID == b.botID {
		return true
	}

	return false
}

// cleanContent はメッセージからボットメンションを削除する
func (b *Bot) cleanContent(content string) string {
	// ボットメンションを削除
	cleaned := strings.Replace(content, "<@"+b.botID+">", "", -1)
	// 余分なスペースを正規化
	return strings.TrimSpace(strings.Join(strings.Fields(cleaned), " "))
}

// processMessage はメッセージを処理して応答を生成する
func (b *Bot) processMessage(m *discordgo.Message) (string, error) {
	// ユーザーメッセージを追加
	userMsg := store.Message{
		Role:      "user",
		Content:   b.cleanContent(m.Content),
		UserID:    m.Author.ID,
		UserName:  m.Author.Username,
		Timestamp: m.Timestamp,
	}

	err := b.store.AddMessage(m.GuildID, m.ChannelID, userMsg)
	if err != nil {
		return "", err
	}

	// 履歴を取得
	history := b.store.GetHistory(m.GuildID, m.ChannelID)

	// AI応答を生成
	response, err := b.aiClient.GenerateResponse(history, b.systemPrompt)
	if err != nil {
		return "", err
	}

	// アシスタントメッセージを追加
	assistantMsg := store.Message{
		Role:      "assistant",
		Content:   response,
		UserID:    b.botID,
		UserName:  "Bot",
		Timestamp: time.Now(),
	}

	err = b.store.AddMessage(m.GuildID, m.ChannelID, assistantMsg)
	if err != nil {
		return "", err
	}

	return response, nil
}

// handleResetCommand はリセットコマンドを処理する
func (b *Bot) handleResetCommand(guildID, channelID string) error {
	return b.store.ClearHistory(guildID, channelID)
}
