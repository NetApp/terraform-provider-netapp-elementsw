# NetApp ElementSW 0.1.0 Example

This repository is designed to demonstrate the capabilities of the [Terraform
NetApp ElementSW Provider][ref-tf-elementsw] at the time of the 0.1.0 release.

[ref-tf-elementsw]: https://www.terraform.io/docs/providers/netapp/elementsw/index.html

This example performs the following:

* Sets up an account. This uses the
  [`elementsw_account` resource][ref-tf-elementsw-account].
* Creates a number of volumes on the cluster tied to the account,
  using the [`elementsw_volume` resource][ref-tf-elementsw-volume].
* Sets up a volume access group for the volumes, using the
  [`elementsw_volume_access_group` resource][ref-tf-elementsw-volume-access-group].
* Finally, creates an initiator tied to the volume access group and volumes using 
  the [`elementsw_initiator` resource][ref-tf-elementsw-initiator].

[ref-tf-elementsw-account]: https://www.terraform.io/docs/providers/netapp/elementsw/r/account.html
[ref-tf-elementsw-initiator]: https://www.terraform.io/docs/providers/netapp/elementsw/r/initiator.html
[ref-tf-elementsw-volume]: https://www.terraform.io/docs/providers/netapp/elementsw/r/volume.html
[ref-tf-elementsw-volume-access-group]: https://www.terraform.io/docs/providers/netapp/elementsw/r/volume_access_group.html

## Requirements

* A working HCI or SolidFire storage cluster.

## Usage Details

You can either clone the entire
[terraform-provider-elementsw][ref-tf-elementsw-github] repository, or download the
`provider.tf`, `variables.tf`, `resources.tf`, and
`terraform.tfvars.example` files into a directory of your choice. Once done,
edit the `terraform.tfvars.example` file, populating the fields with the
relevant values, and then rename it to `terraform.tfvars`. Don't forget to
configure your endpoint and credentials by either adding them to the
`provider.tf` file, or by using enviornment variables. See
[here][ref-tf-elementsw-provider-settings] for a reference on provider-level
configuration values.

[ref-tf-elementsw-github]: https://github.com/terraform-providers/terraform-provider-netapp-elementsw
[ref-tf-elementsw-provider-settings]: https://www.terraform.io/docs/providers/netapp/elementsw/index.html#argument-reference

Once done, run `terraform init`, and `terraform plan` to review the plan, then
`terraform apply` to execute. If you use Terraform 0.11.0 or higher, you can
skip `terraform plan` as `terraform apply` will now perform the plan for you and
ask you confirm the changes.
