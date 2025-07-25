resource "google_compute_network" "vpc" {
  name                    = "chatbot-${var.environment}-network"
  project                 = var.project_id
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name          = "chatbot-${var.environment}-subnet"
  project       = var.project_id
  network       = google_compute_network.vpc.self_link
  region        = var.region
  ip_cidr_range = "10.0.1.0/24"
}

resource "google_compute_firewall" "allow_egress" {
  name    = "chatbot-${var.environment}-allow-egress"
  project = var.project_id
  network = google_compute_network.vpc.name

  direction = "EGRESS"
  priority  = 1000

  allow {
    protocol = "tcp"
  }

  allow {
    protocol = "udp"
  }

  allow {
    protocol = "icmp"
  }

  destination_ranges = ["0.0.0.0/0"]
}
