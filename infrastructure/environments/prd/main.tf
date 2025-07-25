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

variable "bucket_name" {
  description = "GCS bucket name for Terraform state"
  type        = string
}

provider "google" {
  project = var.project_id
  region  = local.region
}

module "api" {
  source = "../../modules/api"

  project_id = var.project_id
}

module "network" {
  source = "../../modules/network"

  project_id  = var.project_id
  environment = var.environment
  region      = local.region
}

module "iam" {
  source = "../../modules/iam"

  project_id  = var.project_id
  environment = var.environment
}

module "secret_manager" {
  source = "../../modules/secret_manager"

  project_id  = var.project_id
  environment = var.environment
}

module "compute" {
  source = "../../modules/compute"

  project_id                    = var.project_id
  environment                   = var.environment
  region                        = local.region
  zone                          = local.zone
  machine_type                  = local.machine_type
  docker_image                  = var.docker_image
  service_account_email         = module.iam.compute_service_account_email
  network_name                  = module.network.network_name
  subnet_name                   = module.network.subnet_name
  discord_bot_token_secret_name = module.secret_manager.discord_bot_token_secret_name
  gemini_api_key_secret_name    = module.secret_manager.gemini_api_key_secret_name

  depends_on = [
    module.api,
    module.network,
    module.iam,
    module.secret_manager
  ]
}

output "secret_setup_instructions" {
  description = "Instructions for setting up secrets"
  value       = module.secret_manager.setup_instructions
}
