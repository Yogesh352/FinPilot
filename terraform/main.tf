terraform {
  required_version = ">= 1.3.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 6.0.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}


resource "google_bigquery_dataset" "finance_db" {
  dataset_id                  = "finance_db"
  location                    = "US"
  delete_contents_on_destroy  = true
  project    = var.project_id
}

resource "google_service_account" "airbyte_sa" {
  account_id   = "airbyte"
  display_name = "Airbyte Service Account"
}

resource "google_service_account" "analytics_sa" {
  account_id   = "analytics"
  display_name = "Analytics/DBT Service Account"
}

resource "google_project_iam_member" "airbyte_bigquery" {
  project = var.project_id
  role    = "roles/bigquery.dataEditor"
  member  = "serviceAccount:${google_service_account.airbyte_sa.email}"
}

resource "google_project_iam_member" "airbyte_jobuser" {
  project = var.project_id
  role    = "roles/bigquery.jobUser"
  member  = "serviceAccount:${google_service_account.airbyte_sa.email}"
}

resource "google_project_iam_member" "analytics_bigquery" {
  project = var.project_id
  role    = "roles/bigquery.dataEditor"
  member  = "serviceAccount:${google_service_account.analytics_sa.email}"
}

resource "google_project_iam_member" "analytics_jobuser" {
  project = var.project_id
  role    = "roles/bigquery.jobUser"
  member  = "serviceAccount:${google_service_account.analytics_sa.email}"
}
