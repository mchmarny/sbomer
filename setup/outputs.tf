# Description: Outputs for the deployment

output "PROVIDER" {
  value       = google_iam_workload_identity_pool_provider.github_provider.name
  description = "Provider ID to use in Auth Actions."
}

output "ACCOUNT" {
  value       = google_service_account.github_actions_user.email
  description = "Service account to use in GitHub Actions."
}

output "BUCKET" {
  value       = google_storage_bucket.report_bucket.name
  description = "GCS Bucket where the resulting files will be placed."
}
