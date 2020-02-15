# Terraform DigitalOcean Generator

This tool generates [Terraform](https://www.terraform.io/) configuration for
[DigitalOcean](https://www.digitalocean.com/).

Terraform is an Infrastructure-as-Code tool to enable easy provisioning and
management of datacentres.

DigitalOcean is a cloud provider which can be controlled using Terraform.

## Building

You will need a current version of [Go](https://golang.org/) installed.

Either Use the supplied `Makefile` to build (`make`), or build using
`go build` (`go build cmd/terraform-do/main.go`).

A copy of the tool will be built in cmd/terraform-do/terraform-do

## Running

Generate a [Personal Access Token](https://cloud.digitalocean.com/account/api/tokens?i=ad3e11)
from the DigitalOcean control panel. Place the resulting token in the
`do_token` environment variable, then run `terraform-do`.

```bash
export do_token="..."
cmd/terraform-do/terraform-do digitalocean.tf
```

This program creates two files:

* `digitalocean.tf`: the main Terraform config
* `digitalocean.tf.imports`: a set of instructions which allows Terraform to
associate the described resources in the Terraform config with the actual
resources in the DigitalOcean account

You should check that the generated output looks correct before using them.
Specifically, anything with an `# UNKNOWN` comment could not be automatically
detected from the DigitalOcean API.

## Using Generated Configs

Be careful! Terraform has a lot of power and if used incorrectly can destroy
all your DigitalOcean resources---test before using the generated resources.

### Initialise Terraform

```bash
terraform init
```

### Import Resources into Terraform

`digitalocean.tf.imports` contains a set of import instructions to be used
with the `terraform import` command.

For example, given a line `digitalocean_domain.jamesam_uk jamesam.uk`, I can
run the following command to tell Terraform that the resource `jamesam_uk`
should be associated with the domain in my DigitalOcean account named
"jamesam.uk".

```bash
terraform import digitalocean_domain.jamesam_uk jamesam.uk`
```

### Check Changes

You can now run the `terraform plan` tool to see what Terraform thinks the
differences are between the config and what it sees at Digitalocean.

If you see any incorrect 'changes', you can edit `digitalocean.tf` to match what
you see and then re-run `terraform plan`.

### Make and Apply Changes

When you edit `digitalocean.tf`, run `terraform plan` to see what changes
Terraform is going to make to your DigitalOcean set up, then run `terraform apply`
to make those changes automatically.

## Implementation Notes

This project was quickly thrown together in half-a-day to meet my personal needs.
As such, there were a few implementation decisions made in the interest of time.

Support for the following digitalocean resources are not included:

* CDN
* Certificate
* Databases
* DNS
* Droplet Snapshots
* Firewall
* Floating IP Assignment
* Load Balancer
* Kubernetes
* Spaces
* Volume Attachment/Snapshot

Notable feature omissions include:

* Links between resources are done via static ID/name fields returned by the DO API
    * finding related terraform resources and using their attributes would have been nicer
* Very little attempt has been made to ensure the TF output is parsable
    * for example, strings are not currently escaped
* HCL is generated using text/template, as the official HCL library does not (yet) support generating HCL
    * in retrospect, TFJSON might have been a more suitable - albeit uglier - option
* Lack of configurability
* No checking for duplicate resource names
* Simplistic parameter processing
* A lot of the code to get resources from DigitalOcean and convert to HCL are a bit same-y
* No verbose mode
