package digitalocean

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/digitalocean/godo"
	"github.com/jamesmichael/terraform-digitalocean/internal"
)

type sshKey struct {
	godo.Key
}

func (p *doProvider) sshKeys(ctx context.Context) ([]td.Resource, error) {
	list := []td.Resource{}

	err := p.paginate(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
		keys, resp, err := p.client.Keys.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		for _, k := range keys {
			list = append(list, &sshKey{k})
		}

		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (k *sshKey) Terraform(w io.Writer) error {
	type tfKey struct {
		TerraformResourceName string
		Name                  string
		PublicKey             string
	}

	rName, err := td.TerraformResourceName(k.Name)
	if err != nil {
		return err
	}
	out := &tfKey{
		TerraformResourceName: rName,
		Name:                  k.Name,
		PublicKey:             k.PublicKey,
	}

	tmpl, err := template.New("key").Parse(`
resource "digitalocean_ssh_key" "{{.TerraformResourceName}}" {
  name       = "{{.Name}}"
  public_key = "{{.PublicKey}}"
}
`)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, out)
}

func (k *sshKey) TerraformImport(w io.Writer) error {
	rName, err := td.TerraformResourceName(k.Name)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "digitalocean_ssh_key.%s %d\n", rName, k.ID)
	return nil
}
