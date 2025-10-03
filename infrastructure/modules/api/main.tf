resource "google_project_service" "required_apis" {
  for_each = toset([
    "compute.googleapis.com",
    "secretmanager.googleapis.com",
    "logging.googleapis.com",
    "monitoring.googleapis.com",
    "iamcredentials.googleapis.com",
  ])

  project = var.project_id
  service = each.value

  disable_dependent_services = false
  disable_on_destroy         = false
}

resource "time_sleep" "wait_for_api_activation" {
  depends_on = [
    google_project_service.required_apis
  ]

  create_duration = "120s"
}
