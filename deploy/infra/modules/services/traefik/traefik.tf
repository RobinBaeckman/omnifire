#------------- k8 namespace -------------#
resource "kubernetes_namespace" "traefik" {
  metadata {
    name = "traefik"
  }
}

#------------- helm -------------#
resource "helm_release" "traefik" {
  repository       = "https://helm.traefik.io/traefik"
  name             = "traefik"
  chart            = "traefik"
  namespace		    = "traefik"
  timeout          = 900
  version		  	 = "20.8.0"
  depends_on = [
    kubernetes_namespace.traefik,
  ]
  values = [
    "${file("${path.cwd}/helm-values/traefik.yaml")}"
  ]
}
