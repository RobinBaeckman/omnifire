#------------- k8 namespace -------------#
resource "kubernetes_namespace" "data-integration" {
  metadata {
	 name = "data-integration"
  }
}

#------------- helm -------------#
resource "helm_release" "kafka" {
  name      = "mykafka"
  chart     = "${path.module}/charts/cp-helm-charts"
  namespace = "data-integration"
  depends_on = [
    kubernetes_namespace.data-integration,
  ]
  values = [
    "${file("${path.cwd}/helm-values/kafka.yaml")}"
  ]
}
