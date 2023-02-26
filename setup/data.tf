# dataset and tables for disco

resource "google_bigquery_dataset" "db" {
  dataset_id    = var.name
  friendly_name = var.name
  description   = "${var.name} dataset"
  location      = "US"
}

// Role binding
resource "google_bigquery_dataset_access" "access" {
  dataset_id    = google_bigquery_dataset.db.dataset_id
  user_by_email = google_service_account.github_actions_user.email
  role          = "OWNER"
}

// Vulnerabilities
resource "google_bigquery_table" "vul_table" {
  dataset_id = google_bigquery_dataset.db.dataset_id
  table_id   = "vul"

  time_partitioning {
    type = "MONTH"
  }

  schema = data.template_file.schema_vul.rendered
}

// Packages
resource "google_bigquery_table" "pkg_table" {
  dataset_id = google_bigquery_dataset.db.dataset_id
  table_id   = "pkg"

  time_partitioning {
    type = "MONTH"
  }

  schema = data.template_file.schema_pkg.rendered
}
