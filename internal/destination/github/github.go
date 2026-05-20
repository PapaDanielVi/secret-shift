package github

import (
	"context"
	"fmt"

	"github.com/PapaDanielVi/secret-shift/internal/destination"
	"github.com/PapaDanielVi/secret-shift/internal/source"
	gogithub "github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

type Destination struct {
	client   *gogithub.Client
	owner    string
	repo     string
	strategy string
}

func New(token, repo, url, strategy string) (*Destination, error) {
	owner, name, err := parseRepo(repo)
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

	return &Destination{
		client:   client,
		owner:    owner,
		repo:     name,
		strategy: strategy,
	}, nil
}

func (d *Destination) Write(ctx context.Context, secrets []source.Secret) error {
	existing, err := d.listExisting(ctx)
	if err != nil {
		return fmt.Errorf("list existing secrets: %w", err)
	}

	for _, s := range secrets {
		_, exists := existing[s.Name]
		if exists && d.strategy == "skip" {
			continue
		}
		if exists && d.strategy == "report" {
			fmt.Printf("  [report] %s %s already exists, skipping\n", s.Type, s.Name)
			continue
		}

		if s.Type == "secret" {
			if err := d.createSecret(ctx, s.Name, s.Value); err != nil {
				return fmt.Errorf("create secret %s: %w", s.Name, err)
			}
		} else {
			if err := d.createVariable(ctx, s.Name, s.Value, exists); err != nil {
				return fmt.Errorf("create variable %s: %w", s.Name, err)
			}
		}
	}

	return nil
}

func (d *Destination) listExisting(ctx context.Context) (map[string]bool, error) {
	existing := make(map[string]bool)

	opts := &gogithub.ListOptions{PerPage: 100}
	for {
		secrets, resp, err := d.client.Actions.ListRepoSecrets(ctx, d.owner, d.repo, opts)
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
		vars, resp, err := d.client.Actions.ListRepoVariables(ctx, d.owner, d.repo, opts)
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

func (d *Destination) createSecret(ctx context.Context, name, value string) error {
	pubKey, _, err := d.client.Actions.GetRepoPublicKey(ctx, d.owner, d.repo)
	if err != nil {
		return fmt.Errorf("get public key: %w", err)
	}

	encrypted, err := encryptSecret(pubKey.GetKey(), value)
	if err != nil {
		return fmt.Errorf("encrypt secret: %w", err)
	}

	_, err = d.client.Actions.CreateOrUpdateRepoSecret(ctx, d.owner, d.repo, &gogithub.EncryptedSecret{
		Name:           name,
		EncryptedValue: encrypted,
		KeyID:          pubKey.GetKeyID(),
	})
	return err
}

func (d *Destination) createVariable(ctx context.Context, name, value string, exists bool) error {
	variable := &gogithub.ActionsVariable{
		Name:  name,
		Value: value,
	}

	if exists {
		_, err := d.client.Actions.UpdateRepoVariable(ctx, d.owner, d.repo, variable)
		return err
	}

	_, err := d.client.Actions.CreateRepoVariable(ctx, d.owner, d.repo, variable)
	return err
}

func parseRepo(repo string) (owner, name string, err error) {
	for i, c := range repo {
		if c == '/' {
			return repo[:i], repo[i+1:], nil
		}
	}
	return "", "", fmt.Errorf("invalid repo format: %s (expected owner/repo)", repo)
}

var _ destination.Destination = (*Destination)(nil)
