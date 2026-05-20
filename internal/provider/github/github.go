package github

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/PapaDanielVi/secret-shift/internal/provider"
	gogithub "github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

const (
	secretTypeSecret = "secret"
	secretTypeEnv    = "env"
	maxPerPage       = 100
)

type Provider struct {
	client   *gogithub.Client
	owner    string
	repo     string
	strategy string
}

func New(token, repoURL, url, strategy string) (*Provider, error) {
	owner, name, err := parseRepo(repoURL)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	var client *gogithub.Client
	if url != "" && url != "https://api.github.com" {
		client, err = gogithub.NewClient(tc).WithEnterpriseURLs(url, url)
		if err != nil {
			return nil, fmt.Errorf("create enterprise github client: %w", err)
		}
	} else {
		client = gogithub.NewClient(tc)
	}

	return &Provider{
		client:   client,
		owner:    owner,
		repo:     name,
		strategy: strategy,
	}, nil
}

func (p *Provider) Read(ctx context.Context) ([]provider.Secret, error) {
	var secrets []provider.Secret

	actionSecrets, err := p.listActionSecrets(ctx)
	if err != nil {
		return nil, fmt.Errorf("list action secrets: %w", err)
	}
	secrets = append(secrets, actionSecrets...)

	variables, err := p.listVariables(ctx)
	if err != nil {
		return nil, fmt.Errorf("list variables: %w", err)
	}
	secrets = append(secrets, variables...)

	return secrets, nil
}

func (p *Provider) Write(ctx context.Context, secrets []provider.Secret) error {
	existing, err := p.listExisting(ctx)
	if err != nil {
		return fmt.Errorf("list existing secrets: %w", err)
	}

	for _, s := range secrets {
		_, exists := existing[s.Name]
		if exists && p.strategy == "skip" {
			continue
		}
		if exists && p.strategy == "report" {
			slog.Info("Secret already exists, skipping", "type", s.Type, "name", s.Name)
			continue
		}

		if s.Type == secretTypeSecret {
			if err := p.createSecret(ctx, s.Name, s.Value); err != nil {
				return fmt.Errorf("create secret %s: %w", s.Name, err)
			}
		} else {
			if err := p.createVariable(ctx, s.Name, s.Value, exists); err != nil {
				return fmt.Errorf("create variable %s: %w", s.Name, err)
			}
		}
	}

	return nil
}

func (p *Provider) listActionSecrets(ctx context.Context) ([]provider.Secret, error) {
	var result []provider.Secret
	opts := &gogithub.ListOptions{PerPage: maxPerPage}

	for {
		secrets, resp, err := p.client.Actions.ListRepoSecrets(ctx, p.owner, p.repo, opts)
		if err != nil {
			return nil, err
		}
		for _, sec := range secrets.Secrets {
			result = append(result, provider.Secret{
				Name: sec.Name,
				Type: secretTypeSecret,
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return result, nil
}

func (p *Provider) listVariables(ctx context.Context) ([]provider.Secret, error) {
	var result []provider.Secret
	opts := &gogithub.ListOptions{PerPage: maxPerPage}

	for {
		vars, resp, err := p.client.Actions.ListRepoVariables(ctx, p.owner, p.repo, opts)
		if err != nil {
			return nil, err
		}
		for _, v := range vars.Variables {
			result = append(result, provider.Secret{
				Name:  v.Name,
				Value: v.Value,
				Type:  secretTypeEnv,
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return result, nil
}

func (p *Provider) listExisting(ctx context.Context) (map[string]bool, error) {
	existing := make(map[string]bool)

	opts := &gogithub.ListOptions{PerPage: maxPerPage}
	for {
		secrets, resp, err := p.client.Actions.ListRepoSecrets(ctx, p.owner, p.repo, opts)
		if err != nil {
			return nil, err
		}
		for _, s := range secrets.Secrets {
			existing[s.Name] = true
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	for {
		vars, resp, err := p.client.Actions.ListRepoVariables(ctx, p.owner, p.repo, opts)
		if err != nil {
			return nil, err
		}
		for _, v := range vars.Variables {
			existing[v.Name] = true
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return existing, nil
}

func (p *Provider) createSecret(ctx context.Context, name, value string) error {
	pubKey, _, err := p.client.Actions.GetRepoPublicKey(ctx, p.owner, p.repo)
	if err != nil {
		return fmt.Errorf("get public key: %w", err)
	}

	encrypted, err := encryptSecret(pubKey.GetKey(), value)
	if err != nil {
		return fmt.Errorf("encrypt secret: %w", err)
	}

	_, err = p.client.Actions.CreateOrUpdateRepoSecret(ctx, p.owner, p.repo, &gogithub.EncryptedSecret{
		Name:           name,
		EncryptedValue: encrypted,
		KeyID:          pubKey.GetKeyID(),
	})
	return err
}

func (p *Provider) createVariable(ctx context.Context, name, value string, exists bool) error {
	variable := &gogithub.ActionsVariable{
		Name:  name,
		Value: value,
	}

	if exists {
		_, err := p.client.Actions.UpdateRepoVariable(ctx, p.owner, p.repo, variable)
		return err
	}

	_, err := p.client.Actions.CreateRepoVariable(ctx, p.owner, p.repo, variable)
	return err
}

func parseRepo(repo string) (string, string, error) {
	for i, c := range repo {
		if c == '/' {
			return repo[:i], repo[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("invalid repo format: %s (expected owner/repo)", repo)
}

var _ provider.Source = (*Provider)(nil)
var _ provider.Destination = (*Provider)(nil)
