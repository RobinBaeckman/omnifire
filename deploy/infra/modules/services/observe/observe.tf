#------------- k8 namespace -------------#
resource "kubernetes_namespace" "observe" {
  metadata {
    name = "observe"
  }
}

#------------- helm -------------#
resource "helm_release" "prometheus-community" {
  atomic           = true
  repository       = "https://prometheus-community.github.io/helm-charts"
  name             = "prometheus"
  chart            = "prometheus"
  namespace		    = "observe"
  timeout          = 900
  version			 = "19.3.3"
  depends_on = [
    kubernetes_namespace.observe,
  ]
  values = [
    "${file("${path.cwd}/helm-values/prometheus.yaml")}"
  ]
}

resource "helm_release" "grafana" {
  atomic           = true
  repository       = "https://grafana.github.io/helm-charts"
  name             = "grafana"
  chart            = "grafana"
  namespace		    = "observe"
  timeout          = 900
  version			 = "6.50.8"
  depends_on = [
    kubernetes_namespace.observe,
  ]
  values = [
    "${file("${path.cwd}/helm-values/grafana.yaml")}"
  ]
}

resource "helm_release" "otel" {
  atomic           = true
  repository       = "https://open-telemetry.github.io/opentelemetry-helm-charts"
  name             = "otel"
  chart            = "opentelemetry-collector"
  namespace		    = "observe"
  timeout          = 900
  version          = "0.47.0"
  depends_on = [
    kubernetes_namespace.observe,
  ]
  values = [
    "${file("${path.cwd}/helm-values/otel.yaml")}"
  ]
}

resource "helm_release" "tempo" {
  atomic           = true
  repository       = "https://grafana.github.io/helm-charts"
  name             = "tempo"
  chart            = "tempo"
  namespace		    = "observe"
  timeout          = 900
  version          = "1.0.0"
  depends_on = [
    kubernetes_namespace.observe,
  ]
  values = [
    "${file("${path.cwd}/helm-values/grafana-tempo.yaml")}"
  ]
}

resource "helm_release" "loki" {
  atomic           = true
  repository       = "https://grafana.github.io/helm-charts"
  name             = "loki"
  chart            = "loki"
  namespace		    = "observe"
  timeout          = 900
  version          = "4.6.1"
  depends_on = [
    kubernetes_namespace.observe,
  ]
  values = [
    "${file("${path.cwd}/helm-values/grafana-loki.yaml")}"
  ]
}

resource "helm_release" "promtail" {
  atomic           = true
  repository       = "https://grafana.github.io/helm-charts"
  name             = "promtail"
  chart            = "promtail"
  namespace		    = "observe"
  timeout          = 900
  version          = "6.8.3"
  depends_on = [
    kubernetes_namespace.observe,
  ]
  values = [
    "${file("${path.cwd}/helm-values/promtail.yaml")}"
  ]
}

resource "helm_release" "pyroscope" {
  atomic           = true
  repository       = "https://pyroscope-io.github.io/helm-chart"
  name             = "pyroscope"
  chart            = "pyroscope"
  namespace		    = "observe"
  timeout          = 900
  version          = "0.2.92"
  depends_on = [
    kubernetes_namespace.observe,
  ]
  values = [
    "${file("${path.cwd}/helm-values/pyroscope.yaml")}"
  ]
}
