terraform {
  required_providers {
    astronomer = {
      source  = "astronomer/astronomer"
      version = ">= 0.1.6"
    }
  }
}

provider "astronomer" {
  api_token = var.astro_api_token
}

resource "astronomer_workspace" "example" {
  label = "my-workspace"
}

resource "astronomer_deployment" "example" {
  workspace_id = astronomer_workspace.example.id
  label        = "my-airflow-deployment"
  region       = "us-central1"
  cloud_provider = "gcp"

  airflow_version = "2.6.3"

  executor = "CeleryExecutor"

  environment_variables = {
    AIRFLOW__CORE__LOAD_EXAMPLES = "False"
  }
}
