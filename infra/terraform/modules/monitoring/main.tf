resource "kubernetes_namespace_v1" "this" {
  metadata {
    name = var.namespace
    labels = {
      "app.kubernetes.io/part-of" = "round-robin"
      "app.kubernetes.io/managed-by" = "terraform"
    }
  }
}

resource "helm_release" "kube_prometheus_stack" {
  name             = "kube-prometheus-stack"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  namespace        = kubernetes_namespace_v1.this.metadata[0].name
  create_namespace = false
  wait             = true

  values = [
    yamlencode({
      grafana = {
        adminPassword = "admin"
        ingress = {
          enabled          = true
          ingressClassName = var.ingress_class_name
          hosts            = [var.grafana_host]
        }
        sidecar = {
          dashboards = {
            enabled         = true
            label           = "grafana_dashboard"
            searchNamespace = "ALL"
          }
        }
      }
      prometheus = {
        ingress = {
          enabled          = true
          ingressClassName = var.ingress_class_name
          hosts            = [var.prometheus_host]
        }
        prometheusSpec = {
          serviceMonitorSelectorNilUsesHelmValues = false
        }
      }
    })
  ]

  depends_on = [
    kubernetes_namespace_v1.this
  ]
}

resource "kubernetes_config_map_v1" "dashboard" {
  metadata {
    name      = "round-robin-overview"
    namespace = kubernetes_namespace_v1.this.metadata[0].name
    labels = {
      grafana_dashboard = "1"
    }
  }

  data = {
    "round-robin-overview.json" = var.dashboard_json
  }

  depends_on = [
    helm_release.kube_prometheus_stack
  ]
}

