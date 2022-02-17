# NetApp ElementSW v20.11 Example

Examples in resources.tf.example1 and resources.tf.example2 are designed to demonstrate the capabilities of the [Terraform
NetApp ElementSW Provider][ref-tf-elementsw].

[ref-tf-elementsw]: https://registry.terraform.io/providers/NetApp/netapp-elementsw/latest

## Requirements

* NetApp HCI, SolidFire or eSDS storage cluster (including Element Demo VM)
* Terraform client

## Getting Started

Clone the Git repository and change directory to the ElementSW examples directory:

```sh
git clone https://github.com/NetApp/terraform-provider-netapp-elementsw
cd terraform-provider-netapp-elementsw/examples/elementsw
```

**NOTE:** Before you continue make sure that volume names, sizes, IQN and other variables from the examples do not conflict with your production environment. Pay special attention when deleting resources because there is no undo. As mentioned in the main README file, you may download Element (SolidFire) Demo VM for safe experimenting.

### Example one: create an account and volumes for CHAP access

To try the first example, `resources.tf.example1`, copy the file to `resources.tf` and examine its contents including tenant and volume names so that you can adjust them if they conflict with your current environment.

Run `terraform init` to doownload NetApp ElementSW Provider.

On `terraform apply`, this example will perform the following:

* Set up a tenant account. This uses the `elementsw_account` resource.
* Creates two volumes for the account using the `elementsw_volume` resource.

`terraform apply` requires certain inputs. You can provide them in `terraform.tfvars` (see `terraform.tfvars.example`) or pass them from the CLI like so:

```sh
terraform apply \
  -var="elementsw_username=admin" \
  -var="elementsw_password=admin" \
  -var="elementsw_cluster=192.168.1.34"
```

On `terraform destroy`, all the resources will be deleted (volumes are purged, not just deleted) without the option to undo. You may need to provide the same variables as above - SolidFire cluster username, password and Management Virtual IP.

After first successful apply, make changes to `resources.tf` and run apply again.

If you want to try the second example, remember to destroy resources with `terraform destroy` and then copy the second example, `resources.tf.example2`, over `resources.tf` (if you had it from the first example). Without a clean-up you may encounter errors due to overlap between resources in the two examples.

### Example two: create an account and volumes VAG access

The second example demonstrates the use of Volume Access Group (VAG) and Initiator resource to creates two additional resources:

* Volume Access Group (VAG) for the volumes, using the `elementsw_volume_access_group` resource.
* Initiator tied to the VAG and the volumes using the `elementsw_initiator` resource.

It also does two things differently from the first example:

* Uses a list of volumes, which is simple but less flexible.
* Lets the SolidFire API to automatically generate tenant secrets - also simple, but less flexible.

Because some variables in this example have values set in `resources.tf` and some have defaults defined in `variables.tf`, the number of variables we have to provide via command line can be less than total number of required variables. For example, `elementsw_username` is already defined in `variables.tf` and `elementsw_initiator` in `resources.tf`, but we can still override the value of former through the CLI.

Like in the first example, check the values of variables in these files and change them to avoid any conflict with existing resources.

```sh
terraform apply \
  -var="elementsw_username=admin" \
  -var="elementsw_password=admin" \
  -var="elementsw_cluster=192.168.1.34" \
  -var="volume_name=testVol" \
  -var="volume_size_list=[1073742000,1073742000]"
```

Note that in this example `volume_size_list` defaults to `[]` (empty list) in order to avoid potential problems. You can change the default value if you want to change this behavior.

To destroy resources just created, run `terraform destroy` (you may need to provide the first three variables).

Descend to examples/elementsw subdirectory and use the sample file with variables to create `terraform.tfvars` and then edit the new file to match your environment:

```sh
cp terraform.tfvars.example terraform.tfvars
vim terraform.tfvars
```

Now run `terraform plan` followed by `terraform apply`. You can omit most variables, but beware of security implications of having `elementsw_password` in plain text file. You may still choose to override certain default variables or variables set in `terraform.tfvars`, especially if they are similar or identical to existing resources.

Destroy with `terraform destroy`, the same way as before.

#### Overriding variable values from the CLI

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
  -var="volume_group_name=testTenant" \
  -var="elementsw_tenant_name=testCluster01"
```

### Add own validation rules

To implement own naming rules or conventions, feel free to create Terraform validation rules.

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

`variables.tf` contains few other example of validation rules (acceptable volume sizes (min 1Gi, max 16TiB), initiator secrets, and volume QoS values).

### Extend

If you wish to extend the scope of this provider with minor features, Terraform [generic provisioners](https://www.terraform.io/docs/language/resources/provisioners/file.html) or vendor provisioners may be a convenient way to achieve that without developing in Go.
