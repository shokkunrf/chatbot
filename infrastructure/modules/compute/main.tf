locals {
  startup_script = templatefile("${path.module}/startup.sh", {
    DOCKER_IMAGE        = var.docker_image
    SECRET_NAME_DISCORD = var.discord_bot_token_secret_name
    SECRET_NAME_GEMINI  = var.gemini_api_key_secret_name
  })
}

resource "google_compute_instance" "chatbot" {
  name         = "chatbot-${var.environment}-instance"
  project      = var.project_id
  machine_type = var.machine_type
  zone         = var.zone

  boot_disk {
    initialize_params {
      image = "cos-cloud/cos-stable"
      size  = 10
      type  = "pd-standard"
    }
  }

  network_interface {
    network    = var.network_name
    subnetwork = var.subnet_name

    access_config {
      # Ephemeral external IP
    }
  }

  metadata = {
    enable-oslogin = "FALSE"
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
