locals {
  environment  = "dev"
  region       = "us-west1"
  zone         = "us-west1-a"
  machine_type = "e2-micro"
}

variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "docker_image" {
  description = "Docker image to run"
  type        = string
}

variable "secret_value_discord_bot_token" {
  description = "Discord bot token"
  type        = string
  sensitive   = true
}

variable "secret_value_gemini_api_key" {
  description = "Gemini API key"
  type        = string
  sensitive   = true
}
