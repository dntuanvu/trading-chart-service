resource "kubernetes_namespace" "trading" {
  metadata {
    name = "trading"
  }
}

resource "kubernetes_deployment" "app" {
  metadata {
    name      = "trading-chart"
    namespace = kubernetes_namespace.trading.metadata[0].name
  }

  spec {
    replicas = 2

    selector {
      match_labels = {
        app = "trading-chart"
      }
    }

    template {
      metadata {
        labels = {
          app = "trading-chart"
        }
      }

      spec {
        container {
          name  = "app"
          image = "dntuanvu/trading-chart-service:latest"
          ports {
            container_port = 50051
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "grpc" {
  metadata {
    name      = "trading-chart-svc"
    namespace = kubernetes_namespace.trading.metadata[0].name
  }

  spec {
    selector = {
      app = "trading-chart"
    }

    port {
      name       = "grpc"
      port       = 50051
      target_port = 50051
    }

    type = "ClusterIP"
  }
}
