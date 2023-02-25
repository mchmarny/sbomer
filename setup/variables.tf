# Description: List of variables which can be provided ar runtime to override the specified defaults 

variable "project_id" {
  description = "GCP Project ID"
  type        = string
  nullable    = false
}

variable "name" {
  description = "Base name to derive everythign else from"
  type        = string
  nullable    = false
  default     = "gcftest"
}

variable "location" {
  description = "Deployment location"
  type        = string
  nullable    = false
  default     = "us-west1"
}

variable "git_repo" {
  description = "GitHub Repo"
  type        = string
  nullable    = false
}
