locals {
  startup_script = templatefile("${path.module}/startup.sh", {
    DOCKER_IMAGE        = var.docker_image
    SECRET_NAME_DISCORD = var.secret_name_discord_bot_token
    SECRET_NAME_GEMINI  = var.secret_name_gemini_api_key
  })
}

resource "google_compute_instance_template" "chatbot" {
  name_prefix  = "chatbot-${var.environment}-"
  project      = var.project_id
  machine_type = var.machine_type

  disk {
    source_image = "cos-cloud/cos-stable"
    disk_size_gb = 10
    disk_type    = "pd-standard"
    auto_delete  = true
    boot         = true
  }

  network_interface {
    network    = var.network_name
    subnetwork = var.subnet_name

    access_config {
      network_tier = "STANDARD"
    }
  }

  metadata = {
    startup-script = local.startup_script
  }

  service_account {
    email  = var.service_account_email
    scopes = ["cloud-platform"]
  }

  tags = ["chatbot-${var.environment}"]

  labels = {
    environment = var.environment
    purpose     = "chatbot"
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "google_compute_instance_group_manager" "chatbot" {
  name               = "chatbot-${var.environment}-mig"
  project            = var.project_id
  zone               = var.zone
  base_instance_name = "chatbot-${var.environment}"
  target_size        = 1

  version {
    instance_template = google_compute_instance_template.chatbot.id
  }
}
