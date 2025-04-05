output "grpc_service_name" {
  value = kubernetes_service.grpc.metadata[0].name
}
