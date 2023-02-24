#module "ingress-controller" {
  #source       = "../../modules/services/ingress-controller"
#}

module "database" {
  source       = "../../modules/services/database"
  helm_timeout = local.helm_timeout
}

module "observe" {
  source = "../../modules/services/observe"
}

module "minio" {
  source = "../../modules/services/minio"
}

#module "data-integration" {
  #source = "../../modules/services/data-integration"
#}
