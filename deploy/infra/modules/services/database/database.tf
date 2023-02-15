#------------- k8 namespace -------------#
resource "kubernetes_namespace" "omnifire" {
  metadata {
    name = "omnifire"
  }
}

#------------- helm -------------#
resource "helm_release" "postgresql" {
  repository       = "https://charts.bitnami.com/bitnami"
  name             = "postgresql"
  chart            = "postgresql"
  namespace		    = "omnifire"
  timeout          = var.helm_timeout
  version			 = "12.1.14"
  depends_on = [
    kubernetes_namespace.omnifire,
  ]
  values = [
    "${file("${path.cwd}/helm-values/postgres.yaml")}"
  ]
}
