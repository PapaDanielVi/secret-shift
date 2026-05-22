<img src="docs/icon.png" alt="secret-shift" style="border-radius: 50%; width: 300px; height: 300px; object-fit: cover;"/>

# SecretShift

[![CI](https://github.com/PapaDanielVi/secret-shift/actions/workflows/test.yml/badge.svg)](https://github.com/PapaDanielVi/secret-shift/actions/workflows/test.yml)
[![Lint](https://github.com/PapaDanielVi/secret-shift/actions/workflows/lint.yml/badge.svg)](https://github.com/PapaDanielVi/secret-shift/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/PapaDanielVi/secret-shift)](https://goreportcard.com/report/github.com/PapaDanielVi/secret-shift)
[![GitHub release](https://img.shields.io/github/v/release/PapaDanielVi/secret-shift)](https://github.com/PapaDanielVi/secret-shift/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/PapaDanielVi/secret-shift.svg)](https://pkg.go.dev/github.com/PapaDanielVi/secret-shift)

SecretShift is a CLI tool for migrating and syncing secrets and environment variables between providers. Move secrets from GitHub, GitLab, HashiCorp Vault, etcd, Kubernetes, or local files to any of those destinations.

## Features

- **7 provider types:** GitHub, GitLab, Vault, etcd, Kubernetes, local file (source + destination)
- **Flexible processing:** Filter by name regex or type, add prefixes/suffixes to secret names
- **Multiple run modes:** One-shot, periodic interval, cron-scheduled, or long-running server
- **Server mode:** HTTP health endpoints (`/healthz`, `/readyz`, `/status`) for Kubernetes deployments
- **Dry-run support:** Preview changes without writing to the destination
- **Encrypted file I/O:** Optional AES-256-GCM encryption for file sources and destinations
- **Conflict handling:** Replace, skip, or report existing secrets at the destination
- **Per-step environment variables:** Separate `SECRET_SHIFT_SRC_*` and `SECRET_SHIFT_DST_*` for source and destination credentials
- **Helm chart:** Deploy to Kubernetes as a Deployment (periodic/server) or CronJob (one-shot)

## Installation

### Homebrew (macOS & Linux)

```bash
brew tap PapaDanielVi/tap
brew install secret-shift
```

### Go Install

```bash
go install github.com/PapaDanielVi/secret-shift@latest
```

### From Source

```bash
go build -o secret-shift .
```

### From Release

Download the binary for your platform from the [releases page](https://github.com/PapaDanielVi/secret-shift/releases).

## Quick Start

### Config File

Create a `secret-shift.json`:

```json
{
  "source": {
    "type": "github",
    "repo": "owner/repo",
    "token": "ghp_xxx"
  },
  "process": {
    "add_prefix": "PROD_",
    "include_regex": "^DB_"
  },
  "destination": {
    "type": "vault",
    "vault_address": "https://vault.example.com",
    "vault_path": "myapp/config",
    "token": "hvs_xxx"
  }
}
```

Run the sync:

```bash
secret-shift sync
```

### Environment Variables

All config keys can be set via environment variables. Source fields use the `SECRET_SHIFT_SRC_` prefix and destination fields use `SECRET_SHIFT_DST_`:

```bash
SECRET_SHIFT_SRC_TYPE=github \
SECRET_SHIFT_SRC_GITHUB_TOKEN=ghp_xxx \
SECRET_SHIFT_DST_TYPE=file \
SECRET_SHIFT_DST_PATH=./secrets.json \
secret-shift sync
```

This per-step convention means you can use different credentials for source and destination, e.g. `SECRET_SHIFT_SRC_GITHUB_TOKEN` and `SECRET_SHIFT_DST_GITHUB_TOKEN`.

## Commands

### `secret-shift sync`

Execute a sync pipeline.

| Flag              | Default               | Description                                              |
| ----------------- | --------------------- | -------------------------------------------------------- |
| `-c, --config`    | `./secret-shift.json` | Path to config file                                      |
| `--periodically`  | `false`               | Run in a loop                                            |
| `--frequency`     | `5m`                  | Interval between syncs (e.g. `1m`, `1h`)                 |
| `--cron`          |                       | Cron expression (e.g. `*/5 * * * *`)                     |
| `--dry-run`       | `false`               | Simulate sync without writing to destination             |
| `--server`        | `false`               | Start HTTP server with health endpoints + periodic sync  |
| `--health-port`   | `8080`                | Port for health HTTP server                              |

### `secret-shift version`

Print the current version.

## Sources

| Source         | What it reads                                           |
| -------------- | ------------------------------------------------------- |
| **github**     | GitHub Actions secrets + environment variables          |
| **gitlab**     | GitLab project-level CI/CD variables                    |
| **vault**      | HashiCorp Vault KV v2 secrets                           |
| **etcd**       | etcd key-value pairs under a prefix                     |
| **kubernetes** | K8s Secrets + ConfigMaps (by name, label, or namespace) |
| **file**       | Local JSON or YAML key-value file (optionally encrypted)|

## Destinations

| Destination    | What it writes                                 | Notes                           |
| -------------- | ---------------------------------------------- | ------------------------------- |
| **file**       | Local JSON or YAML file                        | Optional AES-256-GCM encryption |
| **github**     | GitHub Actions secrets + environment variables | RSA-OAEP encrypted              |
| **gitlab**     | GitLab project-level CI/CD variables           |                                 |
| **vault**      | HashiCorp Vault KV v2                          | Single KV entry                 |
| **etcd**       | etcd key-value store                           | One key per secret              |
| **kubernetes** | K8s Secrets + ConfigMaps                       | Routes by secret type           |

## Server Mode

When running with `--server`, secret-shift starts an HTTP server that exposes:

| Endpoint     | Description                                          |
| ------------ | ---------------------------------------------------- |
| `GET /healthz` | Liveness probe — always returns `{"status":"ok"}` |
| `GET /readyz` | Readiness probe — returns `ready` after first successful sync, `not_ready` if no sync has completed or the last sync failed |
| `GET /status` | Detailed status including sync count, error count, last sync time |

In server mode, syncs run periodically in the background (default every 5 minutes).

## Enterprise and Self-Hosted Providers

### GitHub Enterprise

Set the `url` field to your GitHub Enterprise API endpoint:

```json
{
  "source": {
    "type": "github",
    "repo": "owner/repo",
    "url": "https://git.mycompany.com/api/v3",
    "token": "ghp_xxx"
  }
}
```

### GitLab Self-Hosted

Set the `url` field to your GitLab instance:

```json
{
  "source": {
    "type": "gitlab",
    "project_id": "123",
    "url": "https://gitlab.mycompany.com",
    "token": "glpat-xxx"
  }
}
```

## Processing

The processor runs between source and destination:

- **Type filtering:** `include_types` / `exclude_types` filter by `"env"` or `"secret"`
- **Name filtering:** `include_regex` / `exclude_regex` filter by secret name
- **Name transformation:** `add_prefix` and `add_suffix` modify secret names

## Conflict Strategies

For destinations that support it (GitHub, GitLab):

| Strategy  | Behavior                             |
| --------- | ------------------------------------ |
| `replace` | Overwrite existing secrets (default) |
| `skip`    | Silently skip existing secrets       |
| `report`  | Print conflict info and skip         |

## Environment Variable Reference

### Per-Step Variables

The config system supports separate environment variables for source and destination:

| Variable Pattern                    | Example                              |
| ----------------------------------- | ------------------------------------ |
| `SECRET_SHIFT_SRC_TYPE`             | Source provider type                 |
| `SECRET_SHIFT_SRC_GITHUB_TOKEN`     | GitHub token for source              |
| `SECRET_SHIFT_SRC_GITLAB_TOKEN`     | GitLab token for source              |
| `SECRET_SHIFT_SRC_VAULT_TOKEN`      | Vault token for source               |
| `SECRET_SHIFT_SRC_PATH`             | File path for file source            |
| `SECRET_SHIFT_SRC_KUBE_NAMESPACE`   | K8s namespace for source             |
| `SECRET_SHIFT_DST_TYPE`             | Destination provider type            |
| `SECRET_SHIFT_DST_GITHUB_TOKEN`     | GitHub token for destination         |
| `SECRET_SHIFT_DST_GITLAB_TOKEN`     | GitLab token for destination         |
| `SECRET_SHIFT_DST_VAULT_TOKEN`      | Vault token for destination          |
| `SECRET_SHIFT_DST_PATH`             | File path for file destination       |
| `SECRET_SHIFT_DST_KUBE_NAMESPACE`   | K8s namespace for destination        |
| `SECRET_SHIFT_DRY_RUN`              | Enable dry-run mode                  |

Token resolution order: direct config value → `token_env` field → `SECRET_SHIFT_SRC/DST_<PROVIDER>_TOKEN` env var.

## Examples

See the [`examples/`](examples/) directory for ready-to-use configurations for every source→destination combination (49 total), each with a `config.json` and `README.md`.

## Helm Chart

A Helm chart is available at [`charts/secret-shift/`](charts/secret-shift/). It supports three modes:

| Mode         | Description                          | Kubernetes Resource |
| ------------ | ------------------------------------ | ------------------- |
| `server`     | HTTP health endpoints + periodic sync | Deployment          |
| `periodic`   | Periodic sync loop                   | Deployment          |
| `one-shot`   | Single sync run on a schedule        | CronJob             |

```bash
helm install secret-shift charts/secret-shift/ --set mode=server
```

## Project Structure

```
secret-shift/
  cmd/                    # CLI commands (root, sync, version)
  internal/
    config/               # Config loading, validation, env resolution
    pipeline/             # Sync pipeline (source → process → destination)
    provider/             # Provider interfaces + registry
      github/             # GitHub source + destination
      gitlab/             # GitLab source + destination
      vault/              # HashiCorp Vault source + destination
      etcd/               # etcd source + destination
      kubernetes/         # Kubernetes source + destination
      file/               # File source + destination (JSON/YAML, encrypted)
    server/               # HTTP health server (/healthz, /readyz, /status)
  charts/secret-shift/    # Helm chart for Kubernetes deployment
  examples/               # 49 src→dst example configurations
```

## Dependencies

- **CLI:** [Cobra](https://github.com/spf13/cobra) + [Viper](https://github.com/spf13/viper)
- **APIs:** go-github, GitLab client, Vault API, etcd client, Kubernetes client-go
- **Scheduling:** [cron](https://github.com/robfig/cron)
- **HTTP Server:** Standard library `net/http`

## License

MIT
