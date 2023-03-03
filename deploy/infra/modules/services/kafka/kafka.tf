#------------- k8 namespace -------------#
resource "kubernetes_namespace" "kafka" {
  metadata {
	 name = "kafka"
  }
}

#------------- helm -------------#
resource "helm_release" "kafka" {
  name      = "mykafka"
  chart     = "${path.module}/charts/cp-helm-charts"
  namespace = "kafka"
  depends_on = [
	 kubernetes_namespace.kafka,
  ]
  values = [
	 "${file("${path.cwd}/helm-values/kafka.yaml")}"
  ]
}
