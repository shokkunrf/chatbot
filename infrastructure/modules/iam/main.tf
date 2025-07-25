resource "google_service_account" "compute_sa" {
  account_id   = "chatbot-${var.environment}-compute"
  display_name = "Chatbot ${var.environment} Compute Service Account"
  project      = var.project_id
}

resource "google_project_iam_member" "compute_sa_roles" {
  for_each = toset([
    "roles/secretmanager.secretAccessor",
    "roles/logging.logWriter",
    "roles/monitoring.metricWriter",
  ])

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.compute_sa.email}"
}
