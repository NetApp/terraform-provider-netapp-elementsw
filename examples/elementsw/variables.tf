variable "elementsw_username" {
  type        = string
  default     = "admin"
  description = "The username of the Element cluster admin."
}

variable "elementsw_password" {
  type        = string
  sensitive   = true
  default     = "admin"
  description = "The password of the Element cluster admin."
}

variable "elementsw_cluster" {
  type        = string
  description = "Management Virtual IP (MVIP) of the Element cluster (IPv4 or FQDN)."
}

variable "elementsw_api_version" {
  type        = string
  default     = "11.7"
  description = "The API version of the Element cluster."
}

variable "elementsw_tenant_name" {
  type        = string
  default     = "test-account"
  description = "The Element tenant name."
}

variable "elementsw_tenant_target_secret" {
  type        = string
  sensitive   = true
  default     = ""
  description = "The Element tenant's CHAP target secret."
  validation {
    condition = (length(var.elementsw_tenant_target_secret) >= 12 && length(var.elementsw_tenant_target_secret) <= 16) || length(var.elementsw_tenant_target_secret) == 0
    # condition     = length([var.elementsw_tenant_target_secret]) >= 1
    error_message = "Tenant target secret must be between 12 to 16 ASCII characters in length."
  }
}

variable "elementsw_tenant_initiator_secret" {
  type        = string
  sensitive   = true
  default     = ""
  description = "The Element tenant's CHAP initiator secret."
  validation {
    condition = (length(var.elementsw_tenant_initiator_secret) >= 12 && length(var.elementsw_tenant_initiator_secret) <= 16) || length(var.elementsw_tenant_initiator_secret) == 0
    # condition     = length([var.elementsw_tenant_initiator_secret]) >= 12
    error_message = "Tenant initiator secret must be between 12 to 16 ASCII characters in length."
  }
}

variable "volume_group_name" {
  type        = string
  default     = "test-vag"
  description = "The volume group name (VAG) of the Element cluster."
}

variable "elementsw_initiator" {
  type = map(string)
  default = {
    name  = "iqn.1998-01.com.netapp:test-terraform-000000"
    alias = "test-iqn-tf0"
  }
  description = "The IQN of the iSCSI client."
}

variable "volume_name" {
  type        = string
  default     = "testVol"
  description = "The Element volume name."
}

variable "volume_size_list" {
  type        = list(number)
  default     = []
  description = "The list of one or more volume sizes in bytes."

  validation {
    condition = alltrue([
      for v in var.volume_size_list : (v >= 1073742000 && v <= 17592190000000) ? true : false
    ])
    error_message = "Supported volume sizes are 1 - 16,384 GiB (1073742000 - 17592190000000 bytes)."
  }
}

variable "qos" {
  type = map(number)
  default = {
    "min"   = 100
    "max"   = 200
    "burst" = 300
  }
  description = "The SolidFire storage performance QoS settings to apply."
  validation {
    condition     = var.qos.min >= 50 && var.qos.min <= 15000 && var.qos.min <= var.qos.max && var.qos.max <= 50000 && var.qos.max <= var.qos.burst && var.qos.burst <= 50000
    error_message = "SolidFire has rules that apply to QoS values; additionally, rules for this cluster limit Max and Burst IOPS to <= 50,000."
  }
}

variable "sectorsize_512e" {
  type        = bool
  default     = true
  description = "Emulate 512 byte volume sector size (required for ESXi 7.0 and earlier) or not (4096 bytes)."
}

