variable "elementsw_username" {
  type = string
}

variable "elementsw_password" {
  type = string
}

variable "elementsw_cluster" {
  type = string
}

variable "elementsw_api_version" {
  type = string
  default = "10.0"
}

variable "volume_name" {
  type = string 
}

variable "volume_size_list" {
  type = list
  default = [
    "1073741824",
    "1073741824"
  ]
}
