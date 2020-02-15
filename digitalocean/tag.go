package digitalocean

import (
	"context"
	"io"
	"text/template"

	"github.com/digitalocean/godo"
	"github.com/jamesmichael/terraform-digitalocean/internal"
)

type tag struct {
	godo.Tag
}

func (p *doProvider) tags(ctx context.Context) ([]td.Resource, error) {
	list := []td.Resource{}

	err := p.paginate(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
		domains, resp, err := p.client.Tags.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		for _, t := range domains {
			list = append(list, &tag{t})
		}

		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (t *tag) Terraform(w io.Writer) error {
	type tfTag struct {
		TerraformResourceName string
		Name                  string
	}

	rName, err := td.TerraformResourceName(t.Name)
	if err != nil {
		return err
	}
	out := &tfTag{
		TerraformResourceName: rName,
		Name:                  t.Name,
	}

	tmpl, err := template.New("tag").Parse(`
resource "digitalocean_tag" "{{.TerraformResourceName}}" {
  name = "{{.Name}}"
}
`)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, out)
}

func (t *tag) TerraformImport(w io.Writer) error {
	return nil
}
