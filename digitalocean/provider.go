package digitalocean

import (
	"context"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"github.com/jamesmichael/terraform-digitalocean/internal"
)

type doProvider struct {
	client *godo.Client
}

type providerOption func(*doProvider) error

func NewProvider(opts ...providerOption) (*doProvider, error) {
	p := &doProvider{}

	for _, o := range opts {
		if err := o(p); err != nil {
			return nil, err
		}
	}

	return p, nil
}

type tokenSource struct {
	token string
}

func (ts *tokenSource) Token() (*oauth2.Token, error) {
	t := &oauth2.Token{
		AccessToken: ts.token,
	}
	return t, nil
}

func WithPersonalAccessToken(token string) providerOption {
	source := &tokenSource{token}
	return func(p *doProvider) error {
		oauthClient := oauth2.NewClient(context.Background(), source)
		p.client = godo.NewClient(oauthClient)
		return nil
	}
}

func (p *doProvider) Resources(ctx context.Context) ([]td.Resource, error) {
	var resources []td.Resource

	funcs := []func(context.Context) ([]td.Resource, error){
		p.projects,
		p.tags,
		p.droplets,
		p.domains,
		p.sshKeys,
		p.floatingIPs,
		p.volumes,
	}
	for _, f := range funcs {
		rs, err := f(ctx)
		if err != nil {
			return nil, err
		}
		resources = append(resources, rs...)
	}

	return resources, nil
}

type pageLoader func(context.Context, *godo.ListOptions) (*godo.Response, error)

func (p *doProvider) paginate(ctx context.Context, loadPage pageLoader) error {
	opt := &godo.ListOptions{}
	for {
		resp, err := loadPage(ctx, opt)
		if err != nil {
			return err
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return nil
}
