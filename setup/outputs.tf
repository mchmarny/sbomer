# Description: Outputs for the deployment

output "PROJECT" {
  value       = var.project_id
  description = "Project ID where the resoruces were created. "
}

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

output "DATASET" {
  value       = google_bigquery_dataset.db.dataset_id
  description = "BigQuery dataset where the resulting data will be inserted."
}
