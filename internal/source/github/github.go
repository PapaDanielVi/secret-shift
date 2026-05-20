package github

import (
	"context"
	"fmt"

	"github.com/PapaDanielVi/secret-shift/internal/source"
	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

type Source struct {
	client *github.Client
	owner  string
	repo   string
}

func New(token, repo, url string) (*Source, error) {
	owner, name, err := parseRepo(repo)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	var client *github.Client
	if url != "" && url != "https://api.github.com" {
		client, err = github.NewClient(tc).WithEnterpriseURLs(url, url)
		if err != nil {
			return nil, fmt.Errorf("create enterprise github client: %w", err)
		}
	} else {
		client = github.NewClient(tc)
	}

	return &Source{
		client: client,
		owner:  owner,
		repo:   name,
	}, nil
}

func (s *Source) Read(ctx context.Context) ([]source.Secret, error) {
	var secrets []source.Secret

	actionSecrets, err := s.listActionSecrets(ctx)
	if err != nil {
		return nil, fmt.Errorf("list action secrets: %w", err)
	}
	secrets = append(secrets, actionSecrets...)

	variables, err := s.listVariables(ctx)
	if err != nil {
		return nil, fmt.Errorf("list variables: %w", err)
	}
	secrets = append(secrets, variables...)

	return secrets, nil
}

func (s *Source) listActionSecrets(ctx context.Context) ([]source.Secret, error) {
	var result []source.Secret
	opts := &github.ListOptions{PerPage: 100}

	for {
		secrets, resp, err := s.client.Actions.ListRepoSecrets(ctx, s.owner, s.repo, opts)
		if err != nil {
			return nil, err
		}
		for _, sec := range secrets.Secrets {
			result = append(result, source.Secret{
				Name: sec.Name,
				Type: "secret",
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return result, nil
}

func (s *Source) listVariables(ctx context.Context) ([]source.Secret, error) {
	var result []source.Secret
	opts := &github.ListOptions{PerPage: 100}

	for {
		vars, resp, err := s.client.Actions.ListRepoVariables(ctx, s.owner, s.repo, opts)
		if err != nil {
			return nil, err
		}
		for _, v := range vars.Variables {
			result = append(result, source.Secret{
				Name:  v.Name,
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

func parseRepo(repo string) (owner, name string, err error) {
	for i, c := range repo {
		if c == '/' {
			return repo[:i], repo[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("invalid repo format: %s (expected owner/repo)", repo)
}
