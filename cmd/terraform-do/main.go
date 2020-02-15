package main

import (
	"context"
	"os"

	"github.com/jamesmichael/terraform-digitalocean/digitalocean"
)

func main() {
	terraformOutput := "digitalocean.tf"
	if len(os.Args) > 1 {
		terraformOutput = os.Args[1]
	}

	config, err := os.Create(terraformOutput)
	if err != nil {
		panic(err)
	}
	defer config.Close()

	imports, err := os.Create(terraformOutput + ".imports")
	if err != nil {
		panic(err)
	}
	defer imports.Close()

	token := os.Getenv("do_token")
	p, err := digitalocean.NewProvider(
		digitalocean.WithPersonalAccessToken(token),
	)
	if err != nil {
		panic(err)
	}

	resources, err := p.Resources(context.Background())
	if err != nil {
		panic(err)
	}

	terraformer := digitalocean.NewTerraformer()
	if err := terraformer.TerraformProvider(config); err != nil {
		panic(err)
	}

	if err := terraformer.TerraformResources(config, resources); err != nil {
		panic(err)
	}

	if err := terraformer.TerraformImports(imports, resources); err != nil {
		panic(err)
	}
}
