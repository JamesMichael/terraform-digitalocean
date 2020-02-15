package digitalocean

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"text/template"

	"github.com/digitalocean/godo"
	"github.com/jamesmichael/terraform-digitalocean/internal"
)

type droplet struct {
	godo.Droplet
}

func (p *doProvider) droplets(ctx context.Context) ([]td.Resource, error) {
	list := []td.Resource{}

	err := p.paginate(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
		droplets, resp, err := p.client.Droplets.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		for _, d := range droplets {
			list = append(list, &droplet{d})
		}

		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (d *droplet) Terraform(w io.Writer) error {
	type tfDroplet struct {
		TerraformResourceName string
		Name                  string
		Image                 string
		Region                string
		Size                  string
		Backups               bool
		Monitoring            bool
		IPv6                  bool
		PrivateNetworking     bool
		Tags                  string
		VolumeIDs             string
		UserData              string
	}

	rName, err := td.TerraformResourceName(d.Name)
	if err != nil {
		return err
	}
	out := &tfDroplet{
		TerraformResourceName: rName,
		Name:                  d.Name,
		Region:                d.Region.Slug,
		Size:                  d.Size.Slug,
	}

	if i := d.Image; i != nil {
		if s := i.Slug; s != "" {
			out.Image = `"` + s + `"`
		} else {
			out.Image = strconv.Itoa(i.ID)
		}
	}

	for _, f := range d.Features {
		switch f {
		case "private_networking":
			out.PrivateNetworking = true
		case "backups":
			out.Backups = true
		case "ipv6":
			out.IPv6 = true
		case "monitoring":
			out.Monitoring = true
		case "virtio":
		default:
			log.Printf("unknown feature '%s' attached to droplet '%s'\n", f, d.Name)
		}
	}

	if len(d.Tags) > 0 {
		var tags []string
		for _, t := range d.Tags {
			tags = append(tags, `"`+t+`"`)
		}
		out.Tags = strings.Join(tags, ", ")
	}

	if len(d.VolumeIDs) > 0 {
		var volumes []string
		for _, v := range d.VolumeIDs {
			volumes = append(volumes, `"`+v+`"`)
		}
		out.VolumeIDs = strings.Join(volumes, ", ")
	}

	tmpl, err := template.New("droplet").Parse(`
resource "digitalocean_droplet" "{{.TerraformResourceName}}" {
  name               = "{{.Name}}"
  image              = {{.Image}}
  region             = "{{.Region}}"
  size               = "{{.Size}}"
  backups            = {{.Backups}}
  monitoring         = {{.Monitoring}}
  ipv6               = {{.IPv6}}
  private_networking = {{.PrivateNetworking}}
  resize_disk        = true # UNKNOWN{{with .Tags}}
  tags               = [ {{.}} ]{{end}}{{with .UserData}}
  user_data          = "{{.}}"{{end}}{{with .VolumeIDs}}
  volume_ids         = [ {{.}} ]{{end}}
}
`)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, out)
}

func (d *droplet) TerraformImport(w io.Writer) error {
	rName, err := td.TerraformResourceName(d.Name)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "digitalocean_droplet.%s %d\n", rName, d.ID)
	return nil
}
