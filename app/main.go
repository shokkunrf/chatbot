package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"chatbot/ai"
	"chatbot/bot"
	"chatbot/config"
	"chatbot/store"
)

func main() {
	// 設定を読み込む
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	// ログレベルを設定
	log.Printf("chatbot を起動しています (ログレベル: %s)", cfg.LogLevel)

	// メモリストアを作成
	memStore := store.NewMemoryStore(cfg.MaxHistory)
	log.Println("メモリストアを初期化しました")

	// Geminiクライアントを作成
	aiClient, err := ai.NewGeminiClient(cfg.GeminiAPIKey, cfg.SystemPrompt)
	if err != nil {
		log.Fatalf("Geminiクライアントの作成に失敗しました: %v", err)
	}
	defer aiClient.Close()
	log.Println("Geminiクライアントを初期化しました")

	// Discordボットを作成
	botConfig := bot.Config{
		Token:        cfg.DiscordToken,
		SystemPrompt: cfg.SystemPrompt,
	}
	discordBot, err := bot.NewBot(botConfig, memStore, aiClient)
	if err != nil {
		log.Fatalf("Discordボットの作成に失敗しました: %v", err)
	}
	log.Println("Discordボットを初期化しました")

	// ボットを起動
	err = discordBot.Start()
	if err != nil {
		log.Fatalf("Discordボットの起動に失敗しました: %v", err)
	}
	log.Println("Discordボットが起動しました")

	// シグナルハンドリング
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("シャットダウンシグナルを受信しました")

	// ボットを停止
	err = discordBot.Stop()
	if err != nil {
		log.Printf("Discordボットの停止中にエラーが発生しました: %v", err)
	}

	log.Println("chatbot を終了しました")
}
