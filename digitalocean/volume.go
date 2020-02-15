package digitalocean

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"github.com/digitalocean/godo"
	"github.com/jamesmichael/terraform-digitalocean/internal"
)

type Volume struct {
	godo.Volume
}

func (p *doProvider) volumes(ctx context.Context) ([]td.Resource, error) {
	list := []td.Resource{}

	err := p.paginate(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
		volumes, resp, err := p.client.Storage.ListVolumes(ctx, &godo.ListVolumeParams{ListOptions: opt})
		if err != nil {
			return nil, err
		}

		for _, v := range volumes {
			list = append(list, &Volume{v})
		}

		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (v *Volume) Terraform(w io.Writer) error {
	type tfVolume struct {
		TerraformResourceName string
		Region                string
		Name                  string
		Size                  int64
		Description           string
		FSType                string
		FSLabel               string
		Tags                  string
	}

	rName, err := td.TerraformResourceName(v.Name)
	if err != nil {
		return err
	}
	out := &tfVolume{
		TerraformResourceName: rName,
		Region:                v.Region.Slug,
		Name:                  v.Name,
		Size:                  v.SizeGigaBytes,
		Description:           v.Description,
		FSType:                v.FilesystemType,
		FSLabel:               v.FilesystemLabel,
	}

	tmpl, err := template.New("Volume").Parse(`
resource "digitalocean_volume" "{{.TerraformResourceName}}" {
  region                   = "{{.Region}}"
  name                     = "{{.Name}}"
  size                     = {{.Size}}{{with .Description}}
  description              = "{{.}}"{{end}}{{with .FSType}}
  initial_filesystem_type  = "{{.}}"{{end}}{{with .FSLabel}}
  initial_filesystem_label = "{{.}}"{{end}}{{with .Tags}}
  tags                     = [ {{.}} ]{{end}}
}
`)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, out)
}

func (v *Volume) TerraformImport(w io.Writer) error {
	rName, err := td.TerraformResourceName(v.Name)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "digitalocean_floating_ip.%s %s\n", rName, v.ID)
	return nil
}
