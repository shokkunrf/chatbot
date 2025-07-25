resource "google_secret_manager_secret" "discord_bot_token" {
  secret_id = "chatbot-${var.environment}-discord-bot-token"
  project   = var.project_id

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret" "gemini_api_key" {
  secret_id = "chatbot-${var.environment}-gemini-api-key"
  project   = var.project_id

  replication {
    auto {}
  }
}
