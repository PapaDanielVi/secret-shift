package gitlab

import (
	"context"
	"fmt"

	"github.com/PapaDanielVi/secret-shift/internal/destination"
	"github.com/PapaDanielVi/secret-shift/internal/source"
	gogitlab "gitlab.com/gitlab-org/api/client-go"
)

type Destination struct {
	client    *gogitlab.Client
	projectID string
	strategy  string
}

func New(token, projectID, url, strategy string) (*Destination, error) {
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

	return &Destination{
		client:    client,
		projectID: projectID,
		strategy:  strategy,
	}, nil
}

func (d *Destination) Write(ctx context.Context, secrets []source.Secret) error {
	existing, err := d.listExisting(ctx)
	if err != nil {
		return fmt.Errorf("list existing variables: %w", err)
	}

	for _, s := range secrets {
		_, exists := existing[s.Name]
		if exists && d.strategy == "skip" {
			continue
		}
		if exists && d.strategy == "report" {
			fmt.Printf("  [report] variable %s already exists, skipping\n", s.Name)
			continue
		}

		if exists {
			_, _, err = d.client.ProjectVariables.UpdateVariable(d.projectID, s.Name, &gogitlab.UpdateProjectVariableOptions{
				Value: &s.Value,
			})
		} else {
			_, _, err = d.client.ProjectVariables.CreateVariable(d.projectID, &gogitlab.CreateProjectVariableOptions{
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

func (d *Destination) listExisting(ctx context.Context) (map[string]bool, error) {
	existing := make(map[string]bool)
	opts := &gogitlab.ListProjectVariablesOptions{
		ListOptions: gogitlab.ListOptions{PerPage: 100},
	}

	for {
		vars, resp, err := d.client.ProjectVariables.ListVariables(d.projectID, opts)
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

var _ destination.Destination = (*Destination)(nil)
