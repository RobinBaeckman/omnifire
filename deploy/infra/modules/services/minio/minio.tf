#------------- k8 namespace -------------#
resource "kubernetes_namespace" "minio" {
  metadata {
	 name = "minio"
  }
}

resource "helm_release" "minio" {
  name      = "minio"
  repository = "https://operator.min.io/"
  chart     = "operator"
  namespace = "minio"
  version			 = "4.5.8"
  depends_on = [
    kubernetes_namespace.minio,
  ]
  values = [
    "${file("${path.cwd}/helm-values/minio.yaml")}"
  ]
}
