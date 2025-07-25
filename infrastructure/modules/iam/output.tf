output "compute_service_account_email" {
  description = "Email of the compute service account"
  value       = google_service_account.compute_sa.email
}
