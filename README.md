# Terraform NetApp ElementSW Provider

This is the repository for the Terraform NetApp ElementSW Provider, which can be used with Terraform to configure resources on NetApp HCI or SolidFire storage clusters.

For general information about Terraform, visit the [official website][tf-website] and the [GitHub project page][tf-github].

[tf-website]: https://terraform.io/
[tf-github]: https://github.com/hashicorp/terraform

This provider plugin was initially developed by the SolidFire team for use with internal projects. The provider plugin was refactored to be published and maintained.

This provider was tested with ElementSW versions ranging from 11.1 up to 12.2.

## Naming Conventions

ElementSW does not require resource names to be unique.  They are considered as 'labels' and resources in ElementSW are uniquely identified by 'ids'.  However these ids are not
user friendly, and as they are generated on the fly, they make it difficult to track resources and automate.

This provider assumes that resource names are unique, and enforces it within its scope. This is not an issue if everything is managed through Terraform, but could raise
conflicts if the rule is not respected outside of Terraform.

## Using the Provider

The current version of this provider requires Terraform 0.12 or higher to run.

You will need to build the provider before being able to use it (see [the section below](#building-the-provider))

Note that you need to run `terraform init` to fetch the provider before deploying.

### Provider Documentation

<TBD> The provider is documented [here][tf-elementsw-docs].

Check the provider documentation for details on entering your connection information and how to get started with writing configuration for NetApp ElementSW resources.

[tf-elementsw-docs](website/docs/index.html.markdown)

### Controlling the provider version

Note that you can also control the provider version. This requires the use of a `provider` block in your Terraform configuration if you have not added one already.

The syntax is as follows:

```hcl
provider "netapp-elementsw" {
  version = "~> 1.1"
  ...
}
```

Version locking uses a pessimistic operator, so this version lock would mean anything within the 1.x namespace, including or after 1.1.0. Read more [here][provider-vc] on provider version control.

[provider-vc]: https://www.terraform.io/docs/configuration/providers.html#provider-versions

## Building The Provider

### Prerequisites

If you wish to work on the provider, you'll first need [Go][go-website] installed on your machine (version 1.14+ is **required** to build with current dependencies). You'll also need to correctly setup a [GOPATH][gopath], as well as adding `$GOPATH/bin` to your `$PATH`.

[go-website]: https://golang.org/
[gopath]: http://golang.org/doc/code.html#GOPATH

The following go packages are required to build the provider:

```sh
go get github.com/fatih/structs
go get github.com/hashicorp/terraform
go get github.com/sirupsen/logrus
go get github.com/x-cray/logrus-prefixed-formatter
```

### Cloning the Project

First, you will want to clone the repository to
`$GOPATH/src/github.com/netapp/terraform-provider-netapp-elementsw`:

```sh
mkdir -p $GOPATH/src/github.com/netapp
cd $GOPATH/src/github.com/netapp
git clone https://github.com/NetApp/terraform-provider-netapp-elementsw.git
```

### Running the Build

After the clone has been completed, you can enter the provider directory and
build the provider.

```sh
cd $GOPATH/src/github.com/netapp/terraform-provider-netapp-elementsw
make build
```

### Installing the Local Plugin

After the build is complete, copy the `terraform-provider-netapp-elementsw` binary into the same path as your `terraform` binary, and re-run `terraform init`.

After this, your project-local `.terraform/plugins/ARCH/lock.json` (where `ARCH` matches the architecture of your machine) file should contain a SHA256 sum that matches the local plugin. Run `shasum -a 256` on the binary to verify the values match.

## Developing the Provider

**NOTE:** Before you start work on a feature, please make sure to check the [issue tracker][gh-issues] and existing [pull requests][gh-prs] to ensure that work is not being duplicated. For further clarification, you can also ask in a new issue.

[gh-issues]: https://github.com/netapp/terraform-provider-netapp-elementsw/issues
[gh-prs]: https://github.com/netapp/terraform-provider-netapp-elementsw/pulls

See [Building the Provider](#building-the-provider) for details on building the provider.

## Testing the Provider

**NOTE:** Testing the NetApp ElementSW provider is currently a complex operation as it requires having an ElementSW endpoint to test against, which should be hosting a standard configuration for a HCI or SolidFire cluster.

### Configuring Environment Variables

Most of the tests in this provider require a comprehensive list of environment variables to run. See the individual `*_test.go` files in the [`elementsw/`](elementsw/) directory for more details. The next section also describes how you can manage a configuration file of the test environment variables.

#### Using the `.tf-elementsw-devrc.mk` file

The [`tf-elementsw-devrc.mk.example`](tf-elementsw-devrc.mk.example) file contains an up-to-date list of environment variables required to run the acceptance tests. Copy this to `$HOME/.tf-elementsw-devrc.mk` and change the permissions to something more secure (ie: `chmod 600 $HOME/.tf-elementsw-devrc.mk`), and configure the variables accordingly.

### Running the Acceptance Tests

After this is done, you can run the acceptance tests by running:

```sh
make testacc
```

If you want to run against a specific set of tests, run `make testacc` with the `TESTARGS` parameter containing the run mask as per below:

```sh
make testacc TESTARGS="-run=TestAccElementSwVolume"
```

This following example would run all of the acceptance tests matching `TestAccElementSwVolume`. Change this for the specific tests you want to run.

## Walkthrough example

### Installing Go and Terraform

```sh
bash
mkdir tf_na_elementsw
cd tf_na_elementsw

# if you want a private installation, use
export GO_INSTALL_DIR=`pwd`/go_install
mkdir $GO_INSTALL_DIR
# otherwise, go recommends to use
export GO_INSTALL_DIR=/usr/local


curl -O https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz
tar -C $GO_INSTALL_DIR -xvf go1.14.1.linux-amd64.tar.gz

export PATH=$PATH:$GO_INSTALL_DIR/go/bin

curl -O https://releases.hashicorp.com/terraform/0.12.24/terraform_0.12.24_linux_amd64.zip
unzip terraform_0.12.24_linux_amd64.zip
mv terraform $GO_INSTALL_DIR/go/bin
```

### Installing dependencies

```sh
# make sure git is installed
which git

export GOPATH=`pwd`
go get github.com/fatih/structs
go get github.com/hashicorp/terraform
go get github.com/sirupsen/logrus
go get github.com/x-cray/logrus-prefixed-formatter
```

Note getting the terraform package also builds and installs terraform in $GOPATH/bin.
The version in go/bin is a stable release.

### Cloning the NetApp provider repository and building the provider

```sh
mkdir -p $GOPATH/src/github.com/netapp
cd $GOPATH/src/github.com/netapp
git clone https://github.com/NetApp/terraform-provider-netapp-elementsw.git
cd terraform-provider-netapp-elementsw
make build
mv $GOPATH/bin/terraform-provider-netapp-elementsw $GO_INSTALL_DIR/go/bin
```

The build step will install the provider in the $GOPATH/bin directory. For Terraform v0.11 and v0.12 you could use it from there, for version v0.13 copy it to `/usr/share/terraform/providers/netapp.com/` (see the provided example).

### Sanity check

```shell
cd examples/elementsw/
terraform init
```

Should do nothing but indicate that `Terraform has been successfully initialized!`
