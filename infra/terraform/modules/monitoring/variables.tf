variable "namespace" {
  description = "Namespace for the monitoring stack."
  type        = string
}

variable "grafana_host" {
  description = "Ingress host for Grafana."
  type        = string
}

variable "prometheus_host" {
  description = "Ingress host for Prometheus."
  type        = string
}

variable "ingress_class_name" {
  description = "Ingress class used for monitoring UIs."
  type        = string
  default     = "nginx"
}

variable "dashboard_json" {
  description = "Grafana dashboard JSON content."
  type        = string
}

