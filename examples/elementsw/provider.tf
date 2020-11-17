provider "netapp-elementsw" {
  username = var.elementsw_username
  password = var.elementsw_password
  elementsw_server = var.elementsw_cluster
  api_version = var.elementsw_api_version
}
