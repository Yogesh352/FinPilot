variable "project_id" {
  description = "GCP Project id"
  type        = string
}

variable "billing_account" {
  description = "The GCP billing account ID"
  type        = string
  default     = null
}

variable "region" {
  description = "The GCP region to use for resources"
  type        = string
  default     = "US"
}
