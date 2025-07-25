output "discord_bot_token_secret_name" {
  description = "Discord bot token secret name"
  value       = google_secret_manager_secret.discord_bot_token.name
}

output "gemini_api_key_secret_name" {
  description = "Gemini API key secret name"
  value       = google_secret_manager_secret.gemini_api_key.name
}

output "setup_instructions" {
  description = "Instructions for setting up secret values"
  value       = <<EOF
To set the secret values, run the following commands:

echo -n "YOUR_DISCORD_BOT_TOKEN" | gcloud secrets versions add ${google_secret_manager_secret.discord_bot_token.secret_id} --data-file=-
echo -n "YOUR_GEMINI_API_KEY" | gcloud secrets versions add ${google_secret_manager_secret.gemini_api_key.secret_id} --data-file=-
EOF
}
