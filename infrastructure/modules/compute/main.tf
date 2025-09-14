locals {
  startup_script = templatefile("${path.module}/startup.sh", {
    DOCKER_IMAGE        = var.docker_image
    SECRET_NAME_DISCORD = var.secret_name_discord_bot_token
    SECRET_NAME_GEMINI  = var.secret_name_gemini_api_key
  })
}

resource "google_compute_instance" "chatbot" {
  name         = "chatbot-${var.environment}-instance"
  project      = var.project_id
  machine_type = var.machine_type
  zone         = var.zone

  shielded_instance_config {
    enable_secure_boot          = true
    enable_vtpm                 = true
    enable_integrity_monitoring = true
  }

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
