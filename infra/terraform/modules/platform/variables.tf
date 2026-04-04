variable "namespace" {
  description = "Namespace for the platform release."
  type        = string
}

variable "chart_path" {
  description = "Path to the local Helm chart."
  type        = string
}

variable "install_ingress" {
  description = "Whether to install ingress-nginx."
  type        = bool
  default     = true
}

variable "ingress_class_name" {
  description = "Ingress class name."
  type        = string
  default     = "nginx"
}

variable "ingress_host" {
  description = "Application ingress hostname."
  type        = string
}

variable "app_image_repository" {
  description = "App image repository."
  type        = string
}

variable "app_image_tag" {
  description = "App image tag."
  type        = string
}

variable "worker_image_repository" {
  description = "Worker image repository."
  type        = string
}

variable "worker_image_tag" {
  description = "Worker image tag."
  type        = string
}

variable "monitoring_enabled" {
  description = "Whether to enable ServiceMonitor resources in the app chart."
  type        = bool
  default     = false
}

