#!/bin/bash
set -e

echo "Starting chatbot initialization..."

# Get secrets from Secret Manager
DISCORD_BOT_TOKEN=$(curl -s -H "Authorization: Bearer $(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token | jq -r .access_token)" \
  "https://secretmanager.googleapis.com/v1/${SECRET_NAME_DISCORD}/versions/latest:access" | jq -r .payload.data | base64 -d)

if [ -z "$DISCORD_BOT_TOKEN" ]; then
  echo "Failed to retrieve Discord bot token"
  shutdown now
  exit 1
fi

GEMINI_API_KEY=$(curl -s -H "Authorization: Bearer $(curl -s -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token | jq -r .access_token)" \
  "https://secretmanager.googleapis.com/v1/${SECRET_NAME_GEMINI}/versions/latest:access" | jq -r .payload.data | base64 -d)

if [ -z "$GEMINI_API_KEY" ]; then
  echo "Failed to retrieve Gemini API key"
  shutdown now
  exit 1
fi

# Run Docker container
docker run -d \
  --name chatbot \
  --restart=always \
  -e DISCORD_BOT_TOKEN="$DISCORD_BOT_TOKEN" \
  -e GEMINI_API_KEY="$GEMINI_API_KEY" \
  --log-driver=gcplogs \
  ${DOCKER_IMAGE}

echo "Chatbot container started successfully"
