package digitalocean

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/digitalocean/godo"
	"github.com/jamesmichael/terraform-digitalocean/internal"
)

type FloatingIP struct {
	godo.FloatingIP
}

func (p *doProvider) floatingIPs(ctx context.Context) ([]td.Resource, error) {
	list := []td.Resource{}

	err := p.paginate(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
		ips, resp, err := p.client.FloatingIPs.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		for _, d := range ips {
			list = append(list, &FloatingIP{d})
		}

		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (ip *FloatingIP) Terraform(w io.Writer) error {
	type tfFloatingIP struct {
		TerraformResourceName string
		Region                string
		DropletID             int
	}

	rName, err := td.TerraformResourceName(ip.URN())
	if err != nil {
		return err
	}
	out := &tfFloatingIP{
		TerraformResourceName: rName,
		Region:                ip.Region.Slug,
		DropletID:             ip.Droplet.ID,
	}

	tmpl, err := template.New("FloatingIP").Parse(`
resource "digitalocean_floating_ip" "{{.TerraformResourceName}}" {
  region     = "{{.Region}}"{{with .DropletID}}
  droplet_id = {{.}}{{end}}
}
`)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, out)
}

func (ip *FloatingIP) TerraformImport(w io.Writer) error {
	rName, err := td.TerraformResourceName(ip.URN())
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "digitalocean_floating_ip.%s %s\n", rName, ip.IP)
	return nil
}
