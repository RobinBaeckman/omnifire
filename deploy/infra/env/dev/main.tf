#module "traefik" {
  #source       = "../../modules/services/traefik"
#}

module "postgres" {
  source       = "../../modules/services/postgres"
  helm_timeout = local.helm_timeout
}

module "observe" {
  source = "../../modules/services/observe"
}

#module "minio" {
  #source = "../../modules/services/minio"
#}

#module "kafka" {
  #source = "../../modules/services/kafka"
#}
