package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v68/github"
)

func setupMockServer(secretNames []string, variables []*github.ActionsVariable) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/owner/repo/actions/secrets":
			resp := github.Secrets{
				TotalCount: len(secretNames),
				Secrets:    make([]*github.Secret, len(secretNames)),
			}
			for i, name := range secretNames {
				resp.Secrets[i] = &github.Secret{Name: name}
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		case "/repos/owner/repo/actions/variables":
			resp := github.ActionsVariables{
				TotalCount: len(variables),
				Variables:  variables,
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "not found: %s", r.URL.Path)
		}
	}))
}

func newTestProvider(serverURL string) *Provider {
	p := &Provider{
		client: github.NewClient(nil),
		owner:  "owner",
		repo:   "repo",
	}
	u, _ := url.Parse(serverURL + "/")
	p.client.BaseURL = u
	return p
}

func TestRead_SecretsAndVariables(t *testing.T) {
	vars := []*github.ActionsVariable{
		{Name: "DEPLOY_ENV", Value: "production"},
		{Name: "LOG_LEVEL", Value: "info"},
	}
	server := setupMockServer([]string{"API_KEY", "DB_PASSWORD"}, vars)
	defer server.Close()

	p := newTestProvider(server.URL)

	secrets, err := p.Read(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(secrets) != 4 {
		t.Fatalf("expected 4 secrets, got %d", len(secrets))
	}

	secretCount := 0
	envCount := 0
	for _, s := range secrets {
		if s.Type == "secret" {
			secretCount++
		}
		if s.Type == "env" {
			envCount++
		}
	}
	if secretCount != 2 {
		t.Errorf("expected 2 secrets, got %d", secretCount)
	}
	if envCount != 2 {
		t.Errorf("expected 2 env vars, got %d", envCount)
	}
}

func TestRead_Empty(t *testing.T) {
	server := setupMockServer(nil, nil)
	defer server.Close()

	p := newTestProvider(server.URL)

	secrets, err := p.Read(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secrets) != 0 {
		t.Fatalf("expected 0 secrets, got %d", len(secrets))
	}
}

func TestRead_VariableValues(t *testing.T) {
	vars := []*github.ActionsVariable{
		{Name: "DEPLOY_ENV", Value: "production"},
	}
	server := setupMockServer(nil, vars)
	defer server.Close()

	p := newTestProvider(server.URL)

	secrets, err := p.Read(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secrets) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(secrets))
	}
	if secrets[0].Name != "DEPLOY_ENV" {
		t.Errorf("expected DEPLOY_ENV, got %s", secrets[0].Name)
	}
	if secrets[0].Value != "production" {
		t.Errorf("expected 'production', got %s", secrets[0].Value)
	}
	if secrets[0].Type != "env" {
		t.Errorf("expected type 'env', got %s", secrets[0].Type)
	}
}

func TestParseRepo(t *testing.T) {
	tests := []struct {
		input       string
		wantOwner   string
		wantRepo    string
		expectError bool
	}{
		{"owner/repo", "owner", "repo", false},
		{"my-org/my-repo", "my-org", "my-repo", false},
		{"invalid", "", "", true},
		{"", "", "", true},
	}

	for _, tt := range tests {
		owner, repo, err := parseRepo(tt.input)
		if tt.expectError {
			if err == nil {
				t.Errorf("parseRepo(%q) expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseRepo(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if owner != tt.wantOwner || repo != tt.wantRepo {
			t.Errorf("parseRepo(%q) = %q, %q; want %q, %q", tt.input, owner, repo, tt.wantOwner, tt.wantRepo)
		}
	}
}
