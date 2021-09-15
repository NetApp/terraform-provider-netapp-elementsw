variable "elementsw_username" {
  type        = string
  default     = "admin"
  description = "The username of the Element cluster admin."
}

variable "elementsw_password" {
  type        = string
  sensitive   = true
  description = "The password of the Element cluster admin."
}

variable "elementsw_server" {
  type        = string
  description = "Management Virtual IP (MVIP) of the Element cluster (IPv4 or FQDN)."
}

variable "elementsw_api_version" {
  type        = string
  default     = "11.7"
  description = "The API version of the Element cluster."
}

variable "elementsw_initiator" {
  type        = string
  default     = "iqn.1998-01.com.netapp:test-terraform-000000"
  description = "The IQN of the iSCSI client."
}

variable "volume_name" {
  type        = string
  description = "The Element volume name."
}

variable "volume_size_list" {
  type        = list(number)
  default     = [ ]
  description = "The list of one or more volume sizes in bytes."

  validation {
    condition     = alltrue([
      for v in var.volume_size_list :(v >= 1073742000 && v <= 17592190000000) ? true : false
    ])
    error_message = "Supported volume sizes are 1 - 16,384 GiB (1073742000 - 17592190000000 bytes)."
  }
}

