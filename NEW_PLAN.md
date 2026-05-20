# SecretShift — Implementation Plan

SecretShift is a CLI tool for migrating/syncing secrets and environment variables between sources (GitHub, GitLab, Vault, etcd, K8s) and destinations (GitHub, GitLab, Vault, etcd, K8s, file).

**Pipeline:** `read -> process -> write`

## v1 Scope

- **Sources:** GitHub, GitLab
- **Destinations:** File (JSON/YAML), GitHub, GitLab
- **Flow:** JSON config file + CLI flags (no TUI)
- **Periodic sync:** In-process loop with `--periodically` and `--frequency`
- **Config:** JSON config file + environment variable overrides
- **Tech stack:** Cobra (CLI), Viper (config/env), official API clients

## Architecture

```
cmd/
  root.go              # Entry point, global flags (-c, --periodically, --frequency)
  sync.go              # Subcommand: execute a sync pipeline

internal/
  config/
    config.go          # Config struct, JSON loading, env var binding, validation

  pipeline/
    pipeline.go        # Orchestrates read -> process -> write
    process.go         # Processing logic (prefix, suffix, regex include/exclude, type filter)

  source/
    source.go          # Source interface
    github/            # GitHub Actions secrets + Variables API
    gitlab/            # GitLab CI/CD variables API

  destination/
    destination.go     # Destination interface
    file/              # Write secrets to JSON/YAML file
    github/            # Write to GitHub
    gitlab/            # Write to GitLab
```

## Key Interfaces

```go
// Source reads secrets from a provider
type Source interface {
    Read(ctx context.Context) ([]Secret, error)
}

// Destination writes secrets to a provider
type Destination interface {
    Write(ctx context.Context, secrets []Secret) error
}

// Secret is the canonical unit passed through the pipeline
type Secret struct {
    Name  string
    Value string
    Type  string // "env" or "secret"
}
```

## Config File Structure

```json
{
  "source": {
    "type": "github",
    "repo": "owner/repo",
    "token_env": "GITHUB_TOKEN",
    "url": "https://api.github.com"
  },
  "process": {
    "add_prefix": "PROD_",
    "add_suffix": "_V2",
    "include_regex": "^DB_",
    "exclude_regex": "^DEBUG_",
    "include_types": ["secret"],
    "exclude_types": ["env"]
  },
  "destination": {
    "type": "file",
    "path": "./output/secrets.json",
    "format": "json",
    "conflict_strategy": "replace"
  }
}
```

## Environment Variables

All config values can be overridden via `SECRET_SHIFT_*` env vars:
- `SECRET_SHIFT_SOURCE_TYPE`
- `SECRET_SHIFT_SOURCE_REPO`
- `SECRET_SHIFT_SOURCE_TOKEN` (preferred over token_env for direct token passing)
- `SECRET_SHIFT_DESTINATION_TYPE`
- `SECRET_SHIFT_DESTINATION_PATH`
- `SECRET_SHIFT_GITHUB_TOKEN`
- `SECRET_SHIFT_GITLAB_TOKEN`

## CLI Usage

```bash
# One-shot sync from config
secret-shift sync -c config.json

# Periodic sync every 10 minutes
secret-shift sync -c config.json --periodically --frequency 10m

# Override config values via flags
secret-shift sync -c config.json --source-repo other/repo --dest-path ./out.yaml
```

## Conflict Strategies (write step)

- `replace` — overwrite existing secrets (default)
- `skip` — skip secrets that already exist
- `report` — print a report of what already exists, then skip

## Implementation Order

1. **Project scaffolding** — `go.mod`, Cobra root command, Viper config binding
2. **Config layer** — JSON config structs, file loading, env var overrides, validation
3. **Secret model + Source/Destination interfaces**
4. **Source: GitHub** — list + read Actions secrets and variables
5. **Source: GitLab** — list + read CI/CD variables
6. **Destination: File** — write secrets to JSON or YAML file
7. **Destination: GitHub** — create/update Actions secrets and variables
8. **Destination: GitLab** — create/update CI/CD variables
9. **Process step** — prefix, suffix, regex include/exclude, type filter
10. **Pipeline orchestration** — wire read -> process -> write
11. **Periodic sync** — in-process loop with `--periodically` and `--frequency`
12. **Integration tests** — mock HTTP servers for GitHub/GitLab APIs

## Out of Scope for v1

- TUI flow (Bubble Tea) — v2
- HashiCorp Vault, etcd, Kubernetes sources/destinations — v2
- Cron expression support — v2
- Encryption at rest for file destination — v2

## Verification

1. `go build ./...` — compiles cleanly
2. `go test ./...` — all unit tests pass
3. Manual: run `secret-shift sync -c test-config.json` against a real GitHub repo, verify secrets appear in output file
4. Manual: run with `--periodically --frequency 1m`, verify it loops correctly, then Ctrl+C to stop
