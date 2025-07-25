variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "region" {
  description = "GCP Region"
  type        = string
}

variable "zone" {
  description = "GCP Zone"
  type        = string
}

variable "machine_type" {
  description = "GCE machine type"
  type        = string
}

variable "docker_image" {
  description = "Docker image to run"
  type        = string
}

variable "service_account_email" {
  description = "Service account email for the instance"
  type        = string
}

variable "network_name" {
  description = "VPC network name"
  type        = string
}

variable "subnet_name" {
  description = "Subnet name"
  type        = string
}

variable "discord_bot_token_secret_name" {
  description = "Discord bot token secret name"
  type        = string
}

variable "gemini_api_key_secret_name" {
  description = "Gemini API key secret name"
  type        = string
}
