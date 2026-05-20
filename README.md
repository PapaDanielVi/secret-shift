# SecretShift

SecretShift is a CLI tool for migrating and syncing secrets and environment variables between providers. Move secrets from GitHub, GitLab, HashiCorp Vault, etcd, or Kubernetes to any of those destinations — or to a local file.

## Features

- **6 source/destination types:** GitHub, GitLab, Vault, etcd, Kubernetes, and local file
- **Flexible processing:** Filter by name regex or type, add prefixes/suffixes to secret names
- **Multiple run modes:** One-shot, periodic interval, or cron-scheduled
- **Interactive TUI:** Terminal wizard for guided setup
- **Encrypted file output:** Optional AES-256-GCM encryption for file destinations
- **Conflict handling:** Replace, skip, or report existing secrets at the destination

## Installation

```bash
go build -o secret-shift .
```

## Quick Start

### Interactive TUI

Launch the terminal wizard to configure and run a sync:

```bash
secret-shift tui
```

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

All config keys can be set via environment variables with the `SECRET_SHIFT_` prefix:

```bash
SECRET_SHIFT_SOURCE_TYPE=github \
SECRET_SHIFT_SOURCE_TOKEN=ghp_xxx \
SECRET_SHIFT_DESTINATION_TYPE=file \
SECRET_SHIFT_DESTINATION_PATH=./secrets.json \
secret-shift sync
```

## Commands

### `secret-shift sync`

Execute a sync pipeline.

| Flag | Default | Description |
|------|---------|-------------|
| `-c, --config` | `./secret-shift.json` | Path to config file |
| `--periodically` | `false` | Run in a loop |
| `--frequency` | `5m` | Interval between syncs (e.g. `1m`, `1h`) |
| `--cron` | | Cron expression (e.g. `*/5 * * * *`) |

### `secret-shift tui`

Interactive terminal wizard that walks through source, processing, and destination configuration, then runs the sync.

## Sources

| Source | What it reads |
|--------|--------------|
| **github** | GitHub Actions secrets + environment variables |
| **gitlab** | GitLab project-level CI/CD variables |
| **vault** | HashiCorp Vault KV v2 secrets |
| **etcd** | etcd key-value pairs under a prefix |
| **kubernetes** | K8s Secrets + ConfigMaps (by name, label, or namespace) |

## Destinations

| Destination | What it writes | Notes |
|-------------|---------------|-------|
| **file** | Local JSON or YAML file | Optional AES-256-GCM encryption |
| **github** | GitHub Actions secrets + environment variables | RSA-OAEP encrypted |
| **gitlab** | GitLab project-level CI/CD variables | |
| **vault** | HashiCorp Vault KV v2 | Single KV entry |
| **etcd** | etcd key-value store | One key per secret |
| **kubernetes** | K8s Secrets + ConfigMaps | Routes by secret type |

## Processing

The processor runs between source and destination:

- **Type filtering:** `include_types` / `exclude_types` filter by `"env"` or `"secret"`
- **Name filtering:** `include_regex` / `exclude_regex` filter by secret name
- **Name transformation:** `add_prefix` and `add_suffix` modify secret names

## Conflict Strategies

For destinations that support it (GitHub, GitLab):

| Strategy | Behavior |
|----------|----------|
| `replace` | Overwrite existing secrets (default) |
| `skip` | Silently skip existing secrets |
| `report` | Print conflict info and skip |

## Project Structure

```
secret-shift/
  cmd/                  # CLI commands (root, sync, tui)
  internal/
    config/             # Config loading and validation
    pipeline/           # Sync pipeline (source → process → destination)
    source/             # Source implementations (github, gitlab, vault, etcd, kubernetes)
    destination/        # Destination implementations (file, github, gitlab, vault, etcd, kubernetes)
    tui/                # Bubble Tea interactive TUI
```

## Dependencies

- **CLI:** [Cobra](https://github.com/spf13/cobra) + [Viper](https://github.com/spf13/viper)
- **TUI:** [Bubble Tea v2](https://github.com/charmbracelet/bubbletea) + [Bubbles v2](https://github.com/charmbracelet/bubbles)
- **APIs:** go-github, GitLab client, Vault API, etcd client, Kubernetes client-go
- **Scheduling:** [cron](https://github.com/robfig/cron)

## License

MIT
