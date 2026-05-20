package gitlab

import (
	"context"
	"fmt"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

type Provider struct {
	client    *gogitlab.Client
	projectID string
	strategy  string
}

func New(token, projectID, url, strategy string) (*Provider, error) {
	var client *gogitlab.Client
	var err error

	if url != "" {
		client, err = gogitlab.NewClient(token, gogitlab.WithBaseURL(url))
	} else {
		client, err = gogitlab.NewClient(token)
	}
	if err != nil {
		return nil, fmt.Errorf("create gitlab client: %w", err)
	}

	return &Provider{
		client:    client,
		projectID: projectID,
		strategy:  strategy,
	}, nil
}

func (p *Provider) Read(ctx context.Context) ([]provider.Secret, error) {
	var result []provider.Secret

	opts := &gogitlab.ListProjectVariablesOptions{
		ListOptions: gogitlab.ListOptions{PerPage: 100},
	}

	for {
		vars, resp, err := p.client.ProjectVariables.ListVariables(p.projectID, opts)
		if err != nil {
			return nil, fmt.Errorf("list project variables: %w", err)
		}
		for _, v := range vars {
			result = append(result, provider.Secret{
				Name:  v.Key,
				Value: v.Value,
				Type:  "env",
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return result, nil
}

func (p *Provider) Write(ctx context.Context, secrets []provider.Secret) error {
	existing, err := p.listExisting(ctx)
	if err != nil {
		return fmt.Errorf("list existing variables: %w", err)
	}

	for _, s := range secrets {
		_, exists := existing[s.Name]
		if exists && p.strategy == "skip" {
			continue
		}
		if exists && p.strategy == "report" {
			fmt.Printf("  [report] variable %s already exists, skipping\n", s.Name)
			continue
		}

		if exists {
			_, _, err = p.client.ProjectVariables.UpdateVariable(p.projectID, s.Name, &gogitlab.UpdateProjectVariableOptions{
				Value: &s.Value,
			})
		} else {
			_, _, err = p.client.ProjectVariables.CreateVariable(p.projectID, &gogitlab.CreateProjectVariableOptions{
				Key:   &s.Name,
				Value: &s.Value,
			})
		}
		if err != nil {
			return fmt.Errorf("write variable %s: %w", s.Name, err)
		}
	}

	return nil
}

func (p *Provider) listExisting(ctx context.Context) (map[string]bool, error) {
	existing := make(map[string]bool)
	opts := &gogitlab.ListProjectVariablesOptions{
		ListOptions: gogitlab.ListOptions{PerPage: 100},
	}

	for {
		vars, resp, err := p.client.ProjectVariables.ListVariables(p.projectID, opts)
		if err != nil {
			return nil, err
		}
		for _, v := range vars {
			existing[v.Key] = true
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return existing, nil
}

var _ provider.Source = (*Provider)(nil)
var _ provider.Destination = (*Provider)(nil)
