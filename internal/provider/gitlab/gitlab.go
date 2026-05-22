package gitlab

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

func init() {
	provider.Register(provider.Registration{
		Name: provider.GitLab,
		SourceFactory: func(ctx context.Context, opts map[string]any) (provider.Source, error) {
			token := getString(opts, "token")
			projectID := getString(opts, "project_id")
			url := getString(opts, "url")
			return New(ctx, token, projectID, url, "replace")
		},
		DestFactory: func(ctx context.Context, opts map[string]any) (provider.Destination, error) {
			token := getString(opts, "token")
			projectID := getString(opts, "project_id")
			url := getString(opts, "url")
			strategy := getString(opts, "strategy")
			if strategy == "" {
				strategy = "replace"
			}
			return New(ctx, token, projectID, url, strategy)
		},
	})
}

func getString(opts map[string]any, key string) string {
	if v, ok := opts[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

const maxPerPage = 100

// Provider implements source.Source and destination.Destination for GitLab.
type Provider struct {
	client    *gogitlab.Client
	projectID string
	strategy  string
}

// New creates a GitLab provider for the given project.
func New(_ context.Context, token, projectID, url, strategy string) (*Provider, error) {
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

// Read fetches all project variables from GitLab.
func (p *Provider) Read(_ context.Context) ([]provider.Secret, error) {
	var result []provider.Secret

	opts := &gogitlab.ListProjectVariablesOptions{
		ListOptions: gogitlab.ListOptions{PerPage: maxPerPage},
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

// Write creates or updates project variables in GitLab.
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
			slog.Info("Variable already exists, skipping", "name", s.Name)
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

// listExisting fetches all existing project variable names.
func (p *Provider) listExisting(_ context.Context) (map[string]bool, error) {
	// The GitLab client manages its own request context; ctx is retained in the
	// signature for interface consistency and future use.
	existing := make(map[string]bool)
	opts := &gogitlab.ListProjectVariablesOptions{
		ListOptions: gogitlab.ListOptions{PerPage: maxPerPage},
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
