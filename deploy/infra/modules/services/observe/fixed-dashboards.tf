#------------- custom/fixed grafana dashboards -------------#
# gnetId: 1621
resource "kubernetes_config_map" "k8-cluster-dashboard" {
  metadata {
    name      = "k8-cluster-dashboard"
    namespace = "observe"
    labels = {
      grafana_dashboard = "dashboard"
    }
  }
  depends_on = [
    kubernetes_namespace.observe,
  ]
  data = {
    "k8-cluster-dashboard.json" = "${file("${path.module}/dashboards/k8-cluster-dashboard.json")}"
  }
}

# gnetId: 741
resource "kubernetes_config_map" "k8-deploy-dashboard" {
  metadata {
    name      = "k8-deploy-dashboard"
    namespace = "observe"
    labels = {
      grafana_dashboard = "dashboard"
    }
  }
  depends_on = [
    kubernetes_namespace.observe,
  ]
  data = {
    "k8-deploy-dashboard.json" = "${file("${path.module}/dashboards/k8-deploy-dashboard.json")}"
  }
}

# gnetId: 747
resource "kubernetes_config_map" "k8-pod-dashboard" {
  metadata {
    name      = "k8-pod-dashboard"
    namespace = "observe"
    labels = {
      grafana_dashboard = "dashboard"
    }
  }
  depends_on = [
    kubernetes_namespace.observe,
  ]
  data = {
    "k8-pod-dashboard.json" = "${file("${path.module}/dashboards/k8-pod-dashboard.json")}"
  }
}

# https://raw.githubusercontent.com/confluentinc/cp-helm-charts/master/grafana-dashboard/confluent-open-source-grafana-dashboard.json
resource "kubernetes_config_map" "kafka-dashboard" {
  metadata {
    name      = "kafka-dashboard"
    namespace = "observe"
    labels = {
      grafana_dashboard = "dashboard"
    }
  }
  depends_on = [
    kubernetes_namespace.observe,
  ]
  data = {
    "kafka-dashboard.json" = "${file("${path.module}/dashboards/kafka-dashboard.json")}"
  }
}
