# Terraform NetApp ElementSW Provider

This is the repository for the Terraform NetApp ElementSW Provider, which one can use
with Terraform to work with NetApp HCI or SolidFire storage clusters.

For general information about Terraform, visit the [official
website][tf-website] and the [GitHub project page][tf-github].

[tf-website]: https://terraform.io/
[tf-github]: https://github.com/hashicorp/terraform

This provider plugin was initially developed by the SolidFire team for use with internal
projects. The provider plugin was refactored to be published and maintained.

This provider was tested with ElementSW 11.1 or better.

# Naming Conventions

ElementSW does not require resource names to be unique.  They are considered as 'labels'
and resources in ElementSW are uniquely identified by 'ids'.  However these ids are not
user friendly, and as they are generated on the fly, they make it difficult to track
resources and automate.

This provider assumes that resource names are unique, and enforces it within its scope.
This is not an issue if everything is managed through Terraform, but could raise
conflicts if the rule is not respected outside of Terraform.

# Using the Provider

The current version of this provider requires Terraform v0.12 or higher to
run.

You will need to build the provider before being able to use it
(see [the section below](#building-the-provider)

Note that you need to run `terraform init` to fetch the provider before
deploying.

## Full Provider Documentation

<TBD> The provider is documented [here][tf-elementsw-docs].
Check the provider documentation for details on
entering your connection information and how to get started with writing
configuration for NetApp ElementSW resources.

[tf-elementsw-docs]: https://www.terraform.io/docs/providers/netapp/elementsw/index.html

### Controlling the provider version

Note that you can also control the provider version. This requires the use of a
`provider` block in your Terraform configuration if you have not added one
already.

The syntax is as follows:

```hcl
provider "elementsw" {
  version = "~> 1.1"
  ...
}
```

Version locking uses a pessimistic operator, so this version lock would mean
anything within the 1.x namespace, including or after 1.1.0. [Read
more][provider-vc] on provider version control.

[provider-vc]: https://www.terraform.io/docs/configuration/providers.html#provider-versions

# Building The Provider

## Cloning the Project

First, you will want to clone the repository to
`$GOPATH/src/github.com/netapp/terraform-provider-netapp-elementsw`:

```sh
mkdir -p $GOPATH/src/github.com/netapp
cd $GOPATH/src/github.com/netapp
git clone git@github.com:netapp/terraform-provider-netapp-elementsw
```

## Running the Build

After the clone has been completed, you can enter the provider directory and
build the provider.

```sh
cd $GOPATH/src/github.com/netapp/terraform-provider-netapp-elementsw
make build
```

## Installing the Local Plugin

After the build is complete, copy the `terraform-provider-netapp-elementsw` binary into
the same path as your `terraform` binary, and re-run `terraform init`.

After this, your project-local `.terraform/plugins/ARCH/lock.json` (where `ARCH`
matches the architecture of your machine) file should contain a SHA256 sum that
matches the local plugin. Run `shasum -a 256` on the binary to verify the values
match.

# Developing the Provider

**NOTE:** Before you start work on a feature, please make sure to check the
[issue tracker][gh-issues] and existing [pull requests][gh-prs] to ensure that
work is not being duplicated. For further clarification, you can also ask in a
new issue.

[gh-issues]: https://github.com/netapp/terraform-provider-netapp-elementsw/issues
[gh-prs]: https://github.com/netapp/terraform-provider-netapp-elementsw/pulls

If you wish to work on the provider, you'll first need [Go][go-website]
installed on your machine (version 1.9+ is **required**). You'll also need to
correctly setup a [GOPATH][gopath], as well as adding `$GOPATH/bin` to your
`$PATH`.

[go-website]: https://golang.org/
[gopath]: http://golang.org/doc/code.html#GOPATH

See [Building the Provider](#building-the-provider) for details on building the provider.

# Testing the Provider

**NOTE:** Testing the NetApp ElementSW provider is currently a complex operation as it
requires having an ElementSW endpoint to test against, which should be hosting a
standard configuration for a HCI or SolidFire cluster.

## Configuring Environment Variables

Most of the tests in this provider require a comprehensive list of environment
variables to run. See the individual `*_test.go` files in the
[`elementsw/`](elementsw/) directory for more details. The next section also
describes how you can manage a configuration file of the test environment
variables.

### Using the `.tf-elementsw-devrc.mk` file

The [`tf-elementsw-devrc.mk.example`](tf-elementsw-devrc.mk.example) file contains
an up-to-date list of environment variables required to run the acceptance
tests. Copy this to `$HOME/.tf-elementsw-devrc.mk` and change the permissions to
something more secure (ie: `chmod 600 $HOME/.tf-elementsw-devrc.mk`), and
configure the variables accordingly.

## Running the Acceptance Tests

After this is done, you can run the acceptance tests by running:

```sh
$ make testacc
```

If you want to run against a specific set of tests, run `make testacc` with the
`TESTARGS` parameter containing the run mask as per below:

```sh
make testacc TESTARGS="-run=TestAccElementSwVolume"
```

This following example would run all of the acceptance tests matching
`TestAccElementSwVolume`. Change this for the specific tests you want to
run.
