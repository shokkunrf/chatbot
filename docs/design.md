# 技術設計書

> この文書は、要件定義書 (docs/requirements.md) を元に作成されます。

## 1. アーキテクチャ概要

### システム構成図

```
┌─────────────────┐
│  Discord Users  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Discord Gateway │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│         Chatbot Application             │
├─────────────────────────────────────────┤
│  ┌─────────────┐    ┌───────────────┐  │
│  │   Discord   │    │  Message      │  │
│  │   Handler   │───▶│  Processor    │  │
│  └─────────────┘    └───────┬───────┘  │
│                             │           │
│  ┌─────────────┐    ┌───────▼───────┐  │
│  │  Memory     │◀───│   Gemini      │  │
│  │  Store      │    │   Client      │  │
│  └─────────────┘    └───────────────┘  │
└─────────────────────────────────────────┘
         │
         ▼
┌─────────────────┐
│  Cloud Logging  │
└─────────────────┘
```

### コンポーネント構成

- **Discord Handler**: Discord イベントの受信とレスポンス送信
- **Message Processor**: メッセージの処理とコマンド実行
- **Memory Store**: 会話履歴のメモリ内管理（map構造）
- **Gemini Client**: Gemini API との通信

## 2. 技術スタック

### 言語・フレームワーク

- 言語: Go 1.24+
- Discord ライブラリ: github.com/bwmarrin/discordgo
- AI ライブラリ: github.com/google/generative-ai-go

### インフラ・ミドルウェア

- 実行環境: Docker コンテナ
- ログ: Cloud Logging（stdout/stderr 経由）
- 環境: GCE e2-micro with Container-Optimized OS

## 3. データ設計

### データモデル

```go
// Message - 会話履歴の1メッセージ
type Message struct {
    Role      string    // "user" or "assistant"
    Content   string
    UserID    string
    UserName  string
    Timestamp time.Time
}

// ConversationHistory - チャンネルごとの会話履歴
type ConversationHistory struct {
    Messages []Message // 最大20件
    mu       sync.RWMutex
}

// MemoryStore - 全体の会話履歴管理
type MemoryStore struct {
    conversations map[string]*ConversationHistory // key: "guildID:channelID"
    mu           sync.RWMutex
}
```

### データフロー

1. Discord からメッセージ受信
2. Memory Store から該当チャンネルの履歴取得
3. Gemini API に履歴付きでリクエスト
4. レスポンスを Discord に送信
5. Memory Store に新しいメッセージを追加

## 4. API 設計

### 内部インターフェース

```go
// MemoryStore インターフェース
type Store interface {
    GetHistory(guildID, channelID string) []Message
    AddMessage(guildID, channelID string, message Message)
    ClearHistory(guildID, channelID string)
}

// Gemini Client インターフェース
type AIClient interface {
    GenerateResponse(messages []Message, systemPrompt string) (string, error)
}
```

### Discord コマンド

| コマンド | 説明                               | 権限       |
| -------- | ---------------------------------- | ---------- |
| /reset   | 現在のチャンネルの会話履歴をクリア | 全ユーザー |

## 5. エラーハンドリング

### エラーコード体系

| エラー種別        | 定型文レスポンス                                                       |
| ----------------- | ---------------------------------------------------------------------- |
| Gemini API エラー | 申し訳ありません。現在応答できません。しばらくしてからお試しください。 |
| その他のエラー    | エラーが発生しました。管理者にお問い合わせください。                   |

### エラー処理フロー

1. エラー発生
2. Cloud Logging にエラー詳細を記録
3. Discord に定型文を送信
4. 処理を継続（クラッシュしない）

## 6. セキュリティ設計

### 認証・認可

- Bot トークン: 環境変数で管理
- Gemini API キー: 環境変数で管理

### データ保護

- 個人情報: ユーザーIDとユーザー名のみ保持（メモリ内）
- 会話履歴: アプリケーション再起動で消去

### 入力検証

- メッセージ長: 2000文字以内（Discord の制限）
- 応答長: 175文字以内に制限

## 7. 環境変数

```bash
# 必須
DISCORD_BOT_TOKEN=    # Discord Bot トークン
GEMINI_API_KEY=       # Gemini API キー

# オプション（デフォルト値あり）
LOG_LEVEL=info        # ログレベル (debug, info, warn, error)
MAX_HISTORY=20        # 保持する会話履歴の最大数
RESPONSE_TIMEOUT=5s   # Gemini API のタイムアウト
```

## 8. デプロイ設計

### Docker イメージ

```dockerfile
FROM golang:1.24-bookworm AS builder
FROM gcr.io/distroless/static-debian12:nonroot
```

### 起動コマンド

```bash
docker run -d \
  --name chatbot \
  -e DISCORD_BOT_TOKEN=$DISCORD_BOT_TOKEN \
  -e GEMINI_API_KEY=$GEMINI_API_KEY \
  --restart unless-stopped \
  chatbot:latest
```

## 9. パフォーマンス最適化

### メモリ管理

- 会話履歴: チャンネルあたり最大20メッセージ
- 定期的なガベージコレクション
- メモリ使用量の監視（runtime.MemStats）

### レスポンス時間

- Gemini API: 5秒タイムアウト
- 同期処理でシンプルに実装

## 10. 監視設計

### ログ

- フォーマット: 構造化ログ（JSON）
- 出力先: stdout/stderr → Cloud Logging
- レベル: INFO（本番）、DEBUG（開発）

### メトリクス

- メモリ使用量
- ゴルーチン数
- 処理メッセージ数（ログから集計）
