package digitalocean

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/digitalocean/godo"
	"github.com/jamesmichael/terraform-digitalocean/internal"
)

type domain struct {
	godo.Domain
}

func (p *doProvider) domains(ctx context.Context) ([]td.Resource, error) {
	list := []td.Resource{}

	err := p.paginate(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
		domains, resp, err := p.client.Domains.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		for _, d := range domains {
			list = append(list, &domain{d})
		}

		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (d *domain) Terraform(w io.Writer) error {
	type tfDomain struct {
		TerraformResourceName string
		Name                  string
		IPAddress             string
	}

	rName, err := td.TerraformResourceName(d.Name)
	if err != nil {
		return err
	}
	out := &tfDomain{
		TerraformResourceName: rName,
		Name:                  d.Name,
	}

	tmpl, err := template.New("domain").Parse(`
resource "digitalocean_domain" "{{.TerraformResourceName}}" {
  name = "{{.Name}}"
}
`)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, out)
}

func (d *domain) TerraformImport(w io.Writer) error {
	rName, err := td.TerraformResourceName(d.Name)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "digitalocean_domain.%s %s\n", rName, d.Name)
	return nil
}
