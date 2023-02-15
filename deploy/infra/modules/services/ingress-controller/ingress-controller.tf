#------------- k8 namespace -------------#
resource "kubernetes_namespace" "ingress-controller" {
  metadata {
    name = "ingress-controller"
  }
}

#------------- helm -------------#
resource "helm_release" "traefik" {
  repository       = "https://helm.traefik.io/traefik"
  name             = "traefik"
  chart            = "traefik"
  namespace		    = "ingress-controller"
  timeout          = 900
  version		  	 = "20.8.0"
  depends_on = [
    kubernetes_namespace.traefik,
    #helm_release.metallb,
  ]
  values = [
    "${file("${path.cwd}/helm-values/traefik.yaml")}"
  ]
}
