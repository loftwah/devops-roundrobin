output "application_url" {
  description = "Ingress URL for the application."
  value       = "http://${var.ingress_host}"
}

output "grafana_url" {
  description = "Ingress URL for Grafana."
  value       = var.enable_monitoring ? "http://${var.grafana_host}" : null
}

output "prometheus_url" {
  description = "Ingress URL for Prometheus."
  value       = var.enable_monitoring ? "http://${var.prometheus_host}" : null
}

output "namespace" {
  description = "Application namespace managed by Terraform."
  value       = module.namespace.name
}

