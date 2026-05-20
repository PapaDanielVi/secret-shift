package gitlab

import (
	"context"
	"fmt"

	"github.com/PapaDanielVi/secret-shift/internal/source"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Source struct {
	client    *gitlab.Client
	projectID string
}

func New(token, projectID, url string) (*Source, error) {
	var client *gitlab.Client
	var err error

	if url != "" {
		client, err = gitlab.NewClient(token, gitlab.WithBaseURL(url))
	} else {
		client, err = gitlab.NewClient(token)
	}
	if err != nil {
		return nil, fmt.Errorf("create gitlab client: %w", err)
	}

	return &Source{
		client:    client,
		projectID: projectID,
	}, nil
}

func (s *Source) Read(ctx context.Context) ([]source.Secret, error) {
	var result []source.Secret

	opts := &gitlab.ListProjectVariablesOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	}

	for {
		vars, resp, err := s.client.ProjectVariables.ListVariables(s.projectID, opts)
		if err != nil {
			return nil, fmt.Errorf("list project variables: %w", err)
		}
		for _, v := range vars {
			result = append(result, source.Secret{
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

var _ source.Source = (*Source)(nil)
