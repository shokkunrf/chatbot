output "secret_name_discord_bot_token" {
  description = "Discord bot token secret name"
  value       = google_secret_manager_secret.discord_bot_token.name
}

output "secret_name_gemini_api_key" {
  description = "Gemini API key secret name"
  value       = google_secret_manager_secret.gemini_api_key.name
}
