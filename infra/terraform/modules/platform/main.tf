resource "helm_release" "ingress_nginx" {
  count            = var.install_ingress ? 1 : 0
  name             = "ingress-nginx"
  repository       = "https://kubernetes.github.io/ingress-nginx"
  chart            = "ingress-nginx"
  namespace        = "ingress-nginx"
  create_namespace = true
  wait             = true

  values = [
    yamlencode({
      controller = {
        kind = "DaemonSet"
        ingressClassResource = {
          default = true
        }
        hostPort = {
          enabled = true
        }
        service = {
          type = "ClusterIP"
        }
        nodeSelector = {
          "ingress-ready"   = "true"
          "kubernetes.io/os" = "linux"
        }
        tolerations = [
          {
            key      = "node-role.kubernetes.io/control-plane"
            operator = "Exists"
            effect   = "NoSchedule"
          },
          {
            key      = "node-role.kubernetes.io/master"
            operator = "Exists"
            effect   = "NoSchedule"
          }
        ]
      }
      defaultBackend = {
        enabled = true
      }
    })
  ]
}

resource "helm_release" "platform" {
  name             = "round-robin"
  chart            = var.chart_path
  namespace        = var.namespace
  create_namespace = false
  wait             = true

  values = [
    yamlencode({
      ingress = {
        enabled   = true
        className = var.ingress_class_name
        host      = var.ingress_host
      }
      app = {
        image = {
          repository = var.app_image_repository
          tag        = var.app_image_tag
          pullPolicy = "IfNotPresent"
        }
      }
      worker = {
        image = {
          repository = var.worker_image_repository
          tag        = var.worker_image_tag
          pullPolicy = "IfNotPresent"
        }
      }
      monitoring = {
        enabled = var.monitoring_enabled
      }
    })
  ]

  depends_on = [
    helm_release.ingress_nginx
  ]
}

