#------------- k8 namespace -------------#
resource "kubernetes_namespace" "monitor" {
  metadata {
    name = "monitor"
  }
}

#------------- helm -------------#
resource "helm_release" "prometheus-community" {
  atomic           = true
  repository       = "https://prometheus-community.github.io/helm-charts"
  name             = "prometheus"
  chart            = "prometheus"
  namespace		    = "monitor"
  timeout          = 900
  version			 = "19.3.3"
  depends_on = [
    kubernetes_namespace.monitor,
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
  namespace		    = "monitor"
  timeout          = 900
  version			 = "6.50.7"
  depends_on = [
    kubernetes_namespace.monitor,
	 helm_release.prometheus-community
  ]
  values = [
    "${file("${path.cwd}/helm-values/grafana.yaml")}"
  ]
}

#------------- custom/fixed grafana dashboards -------------#
# gnetId: 1621
resource "kubernetes_config_map" "k8-cluster-dashboard" {
  metadata {
    name      = "k8-cluster-dashboard"
    namespace = "monitor"
    labels = {
      grafana_dashboard = "dashboard"
    }
  }
  data = {
    "k8-cluster-dashboard.json" = "${file("${path.module}/dashboards/k8-cluster-dashboard.json")}"
  }
}

# gnetId: 741
resource "kubernetes_config_map" "k8-deploy-dashboard" {
  metadata {
    name      = "k8-deploy-dashboard"
    namespace = "monitor"
    labels = {
      grafana_dashboard = "dashboard"
    }
  }
  data = {
    "k8-deploy-dashboard.json" = "${file("${path.module}/dashboards/k8-deploy-dashboard.json")}"
  }
}

# gnetId: 747
resource "kubernetes_config_map" "k8-pod-dashboard" {
  metadata {
    name      = "k8-pod-dashboard"
    namespace = "monitor"
    labels = {
      grafana_dashboard = "dashboard"
    }
  }
  data = {
    "k8-pod-dashboard.json" = "${file("${path.module}/dashboards/k8-pod-dashboard.json")}"
  }
}

# https://raw.githubusercontent.com/confluentinc/cp-helm-charts/master/grafana-dashboard/confluent-open-source-grafana-dashboard.json
resource "kubernetes_config_map" "kafka-dashboard" {
  metadata {
    name      = "kafka-dashboard"
    namespace = "monitor"
    labels = {
      grafana_dashboard = "dashboard"
    }
  }
  data = {
    "kafka-dashboard.json" = "${file("${path.module}/dashboards/kafka-dashboard.json")}"
  }
}
