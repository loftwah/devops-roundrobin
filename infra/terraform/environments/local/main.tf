module "namespace" {
  source = "../../modules/namespace"

  name = var.namespace
  labels = {
    "app.kubernetes.io/part-of"    = "round-robin"
    "app.kubernetes.io/managed-by" = "terraform"
  }
}

module "monitoring" {
  count  = var.enable_monitoring ? 1 : 0
  source = "../../modules/monitoring"

  namespace          = var.monitoring_namespace
  grafana_host       = var.grafana_host
  prometheus_host    = var.prometheus_host
  ingress_class_name = var.ingress_class_name
  dashboard_json     = file("${path.root}/../../../../dashboards/round-robin-overview.json")
}

module "platform" {
  source = "../../modules/platform"

  namespace               = module.namespace.name
  chart_path              = "${path.root}/../../../helm/round-robin"
  install_ingress         = var.install_ingress
  ingress_class_name      = var.ingress_class_name
  ingress_host            = var.ingress_host
  app_image_repository    = var.app_image_repository
  app_image_tag           = var.app_image_tag
  worker_image_repository = var.worker_image_repository
  worker_image_tag        = var.worker_image_tag
  monitoring_enabled      = var.enable_monitoring

  depends_on = [
    module.namespace,
    module.monitoring
  ]
}
