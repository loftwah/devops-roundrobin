variable "kubeconfig_path" {
  description = "Path to the kubeconfig used for the local cluster."
  type        = string
  default     = "~/.kube/config"
}

variable "kube_context" {
  description = "Kubernetes context for the local cluster."
  type        = string
  default     = "kind-round-robin"
}

variable "namespace" {
  description = "Namespace for the application release."
  type        = string
  default     = "round-robin"
}

variable "monitoring_namespace" {
  description = "Namespace for the monitoring stack."
  type        = string
  default     = "monitoring"
}

variable "install_ingress" {
  description = "Whether Terraform should install ingress-nginx."
  type        = bool
  default     = false
}

variable "enable_monitoring" {
  description = "Whether Terraform should install kube-prometheus-stack and wire ServiceMonitors."
  type        = bool
  default     = false
}

variable "ingress_class_name" {
  description = "IngressClass name used by the local cluster."
  type        = string
  default     = "nginx"
}

variable "app_image_repository" {
  description = "Container image repository for the app."
  type        = string
  default     = "round-robin-app"
}

variable "app_image_tag" {
  description = "Container image tag for the app."
  type        = string
  default     = "local"
}

variable "worker_image_repository" {
  description = "Container image repository for the worker."
  type        = string
  default     = "round-robin-worker"
}

variable "worker_image_tag" {
  description = "Container image tag for the worker."
  type        = string
  default     = "local"
}

variable "ingress_host" {
  description = "Hostname for the application ingress."
  type        = string
  default     = "app.127.0.0.1.nip.io"
}

variable "grafana_host" {
  description = "Hostname for the Grafana ingress."
  type        = string
  default     = "grafana.127.0.0.1.nip.io"
}

variable "prometheus_host" {
  description = "Hostname for the Prometheus ingress."
  type        = string
  default     = "prometheus.127.0.0.1.nip.io"
}
