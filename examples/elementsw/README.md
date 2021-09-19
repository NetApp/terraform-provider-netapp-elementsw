# NetApp ElementSW v20.11 Example

This repository is designed to demonstrate the capabilities of the [Terraform
NetApp ElementSW Provider][ref-tf-elementsw].

[ref-tf-elementsw]: https://registry.terraform.io/providers/NetApp/netapp-elementsw/latest

On `terraform apply`, this example performs the following:

* Sets up an account. This uses the `elementsw_account` resource
* Creates a number of volumes on the cluster tied to the account,
  using the `elementsw_volume` resource.
* Sets up a volume access group (VAG) for the volumes, using the
  `elementsw_volume_access_group` resource.
* Finally, creates an initiator tied to the volume access group and volumes using
  the `elementsw_initiator` resource.

On `terraform destroy`, it removes all the resources (volumes are purged, not just deleted).

## Requirements

* NetApp HCI, SolidFire or eSDS storage cluster (including Element Demo VM)
* Terraform client

## Getting Started

Clone the Git repository and change directory to the ElementSW examples directory:

```sh
git clone https://github.com/NetApp/terraform-provider-netapp-elementsw
cd terraform-provider-netapp-elementsw/examples/elementsw
terraform init
```

`terraform init` should download ElementSW provider to `.terraform` in the current path.

**NOTE:** Before you continue make sure that volume names, sizes, IQN and other variables from the examples do not conflict with your production environment. Pay special attention when referencing volumes when volumes get deleted earlier than others (so that `testVol-5` becomes one of only three volumes you're working with, and you need to delete the second volume named `testVol-3`). As mentioned in the main README file, you may download Element (SolidFire) Demo VM for safe experimenting.

What's in the files?

- provider.tf - provider settings
- resources.tf - description of resources the provider can manage
- variables.tf - description of some variables that may be used with teh provider (also see terraform.tfvars.example)
- version.tf - provider version control file that downloads ElementSW provider from Terraform Registry. To load own copy, please see the main README file (about building from source)
- terraform.tfvars.example - sample file with variables

### Without terraform.tfvars

Without a variables file we need to make sure Terraform has all of the required variables.

Because some variables in this example have values set in `resources.tf` and some have defaults defined in `variables.tf`, the number of variables we have to provide via command line can be less than the total number of required variables. For example, `elementsw_username` is already defined in `variables.tf` and `elementsw_initiator` in `resources.tf`, but we can still override the value of former through the CLI.

```sh
terraform apply \
  -var="elementsw_username=admin" \
  -var="elementsw_password=admin" \
  -var="elementsw_cluster=192.168.1.34" \
  -var="volume_name=testVol" \
  -var="volume_size_list=[1073742000,1073742000]"
```

Note that in this example `volume_size_list` defaults to `[]` (empty list) in  to avoid potential problems when testing. You can change the default value if you want to change this behavior.

To destroy resources just created, run `terraform destroy` (you need to provide the same set of variables).

### With terraform.tfvars

Descend to examples/elementsw subdirectory and use the sample file to create `terraform.tfvars` and then edit the new file to match your environment:

```sh
cp terraform.tfvars.example terraform.tfvars
vim terraform.tfvars
```

Now run `terraform plan` followed by `terraform apply`. You may still choose to override certain default variables or variables set in `terraform.tfvars`.

Destroy with `terraform destroy`, the same as in the first example.

### Overriding map values from the CLI

This example only shows how values for two maps (QoS and IQN) can be provided from the CLI (Bash shell on Linux). Variations of this approach may be required for different OS.

```sh
terraform apply \
  -var="elementsw_username=admin" \
  -var="elementsw_password=admin" \
  -var="elementsw_cluster=192.168.1.34" \
  -var="volume_name=testVol" \
  -var="volume_size_list=[1073742000,1073742000,1073742000]" \
  -var="sectorsize_512e=false" \
  -var="qos={min=100,max=200,burst=300}" \
  -var="volume_name=dc1-testVol-master" \
  -var="elementsw_initiator={name=\"iqn.1998-01.com.vmware:test-cluster-000001\",alias=\"testNode1\"}" \
  -var="voume_group_name=testTenant" \
  -var="elementsw_tenant_name=testCluster01"
```

### Add own validation rules

To implement own naming rules or conventions, feel free to add Terraform validation rules.

In this example we want to ensure that volume names begin with `dc1`.

```hcl
variable "volume_name" {
  type        = string
  description = "The Element volume name."

  validation {
    condition     = length(var.volume_name) > 2 && substr(var.volume_name, 0, 3) == "dc1"
    error_message = "The volume name string must begin with \"dc1\" and have 3 or more characters."
  }
}
```

Another useful example is a validation rule for acceptable volume sizes (min 1Gi, max 16TiB) - see `variables.tf`.

### Extend

If you wish to extend the scope of this provider with minor features, Terraform [generic provisioners](https://www.terraform.io/docs/language/resources/provisioners/file.html) or vendor provisioners may be a convenient way to achieve that without developing in Go.
