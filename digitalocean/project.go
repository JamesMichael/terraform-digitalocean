package digitalocean

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/digitalocean/godo"
	"github.com/jamesmichael/terraform-digitalocean/internal"
)

type doProject struct {
	godo.Project
	resources []godo.ProjectResource
}

func (p *doProvider) projects(ctx context.Context) ([]td.Resource, error) {
	list := []td.Resource{}

	err := p.paginate(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
		projects, resp, err := p.client.Projects.List(ctx, opt)
		if err != nil {
			return nil, err
		}

		for _, project := range projects {
			resources, err := p.projectResources(ctx, project)
			if err != nil {
				return nil, err
			}
			dop := &doProject{
				Project:   project,
				resources: resources,
			}

			list = append(list, dop)
		}
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (p *doProvider) projectResources(ctx context.Context, project godo.Project) ([]godo.ProjectResource, error) {
	var list []godo.ProjectResource

	err := p.paginate(ctx, func(ctx context.Context, opt *godo.ListOptions) (*godo.Response, error) {
		resources, resp, err := p.client.Projects.ListResources(ctx, project.ID, opt)
		if err != nil {
			return nil, err
		}

		list = append(list, resources...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (p *doProject) Terraform(w io.Writer) error {
	type tfProject struct {
		TerraformResourceName string
		Name                  string
		Description           string
		Purpose               string
		Environment           string
		Resources             string
	}

	rName, err := td.TerraformResourceName(p.Name)
	if err != nil {
		return err
	}
	out := &tfProject{
		TerraformResourceName: rName,
		Name:                  p.Name,
		Description:           p.Description,
		Purpose:               p.Purpose,
		Environment:           p.Environment,
	}

	if len(p.resources) > 0 {
		var resources []string
		for _, r := range p.resources {
			resources = append(resources, `"`+r.URN+`"`)
		}
		out.Resources = strings.Join(resources, ", ")
	}

	tmpl, err := template.New("project").Parse(`
resource "digitalocean_project" "{{.TerraformResourceName}}" {
  name        = "{{.Name}}"
  description = "{{.Description}}"{{with .Purpose}}
  purpose     = "{{.}}"{{end}}{{with .Environment}}
  environment = "{{.}}"{{end}}
  resources   = [ {{.Resources}} ]
}
`)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, out)
}

func (p *doProject) TerraformImport(w io.Writer) error {
	rName, err := td.TerraformResourceName(p.Name)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "digitalocean_project.%s %s\n", rName, p.ID)
	return nil
}
