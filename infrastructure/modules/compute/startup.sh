#!/bin/bash

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

# Stop and remove existing container if it exists
docker stop chatbot 2>/dev/null || true
docker rm chatbot 2>/dev/null || true

# Pull latest Docker image
echo "Pulling Docker image..."
docker pull ${DOCKER_IMAGE} > /dev/null

if [ $? -ne 0 ]; then
  echo "Failed to pull Docker image"
  shutdown now
  exit 1
fi

# Run Docker container
echo "Starting Docker container..."
docker run -d \
  --name chatbot \
  -e DISCORD_BOT_TOKEN="$DISCORD_BOT_TOKEN" \
  -e GEMINI_API_KEY="$GEMINI_API_KEY" \
  --log-driver=gcplogs \
  ${DOCKER_IMAGE} > /dev/null

if [ $? -ne 0 ]; then
  echo "Failed to start Docker container"
  shutdown now
  exit 1
fi

# Clean up unused Docker images (running container's image will be protected)
echo "Cleaning up unused Docker images..."
docker image prune -a -f > /dev/null 2>&1 || true

echo "Chatbot container started successfully"
