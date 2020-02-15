package digitalocean

import (
	"fmt"
	"io"
	"reflect"

	"github.com/jamesmichael/terraform-digitalocean/internal"
)

type terraformer struct {
}

func NewTerraformer() td.Terraformer {
	return &terraformer{}
}

func (tf *terraformer) TerraformImport(w io.Writer, r td.Resource) error {
	type terraformable interface {
		TerraformImport(io.Writer) error
	}
	if t, ok := r.(terraformable); ok {
		return t.TerraformImport(w)
	}
	return nil
}

func (tf *terraformer) TerraformImports(w io.Writer, rs []td.Resource) error {
	for _, r := range rs {
		if err := tf.TerraformImport(w, r); err != nil {
			return err
		}
	}
	return nil
}

func (tf *terraformer) TerraformProvider(w io.Writer) error {
	fmt.Fprint(w, `variable "do_token" {}

provider "digitalocean" {
  token = var.do_token
}
`)

	return nil
}

func (tf *terraformer) TerraformResource(w io.Writer, r td.Resource) error {
	type terraformable interface {
		Terraform(io.Writer) error
	}
	if t, ok := r.(terraformable); ok {
		return t.Terraform(w)
	}

	return fmt.Errorf("resource '%s' does not implement Terraform method", reflect.TypeOf(r))
}

func (tf *terraformer) TerraformResources(w io.Writer, rs []td.Resource) error {
	for _, r := range rs {
		if err := tf.TerraformResource(w, r); err != nil {
			return err
		}
	}
	return nil
}
