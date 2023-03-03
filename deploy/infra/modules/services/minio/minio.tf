#------------- k8 namespace -------------#
resource "kubernetes_namespace" "minio" {
  metadata {
	 name = "minio"
  }
}

resource "helm_release" "minio" {
  name      = "minio"
  repository = "https://charts.min.io/"
  chart     = "minio"
  namespace = "minio"
  version			 = "5.0.7"
  depends_on = [
    kubernetes_namespace.minio,
  ]
  values = [
    "${file("${path.cwd}/helm-values/minio.yaml")}"
  ]
}
