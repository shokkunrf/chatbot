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

  project_id                     = var.project_id
  environment                    = var.environment
  secret_value_discord_bot_token = var.secret_value_discord_bot_token
  secret_value_gemini_api_key    = var.secret_value_gemini_api_key
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
  secret_name_discord_bot_token = module.secret_manager.secret_name_discord_bot_token
  secret_name_gemini_api_key    = module.secret_manager.secret_name_gemini_api_key

  depends_on = [
    module.api,
    module.network,
    module.iam,
    module.secret_manager
  ]
}
