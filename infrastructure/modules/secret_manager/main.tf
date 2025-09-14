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

resource "null_resource" "set_secrets" {
  provisioner "local-exec" {
    command = "${path.module}/set_secrets.sh"
    environment = {
      SECRET_VALUE_DISCORD = var.secret_value_discord_bot_token
      SECRET_VALUE_GEMINI  = var.secret_value_gemini_api_key
      SECRET_NAME_DISCORD  = google_secret_manager_secret.discord_bot_token.secret_id
      SECRET_NAME_GEMINI   = google_secret_manager_secret.gemini_api_key.secret_id
    }
  }

  depends_on = [
    google_secret_manager_secret.discord_bot_token,
    google_secret_manager_secret.gemini_api_key
  ]

  triggers = {
    secret_hash_discord = sha256(var.secret_value_discord_bot_token)
    secret_hash_gemini  = sha256(var.secret_value_gemini_api_key)
  }
}
