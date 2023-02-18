#module "ingress-controller" {
  #source       = "../../modules/services/ingress-controller"
#}

#module "database" {
  #source       = "../../modules/services/database"
  #helm_timeout = local.helm_timeout
#}

module "monitoring" {
  source = "../../modules/services/monitoring"
}

#module "data-integration" {
  #source = "../../modules/services/data-integration"
#}
