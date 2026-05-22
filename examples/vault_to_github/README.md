# Vault to GitHub

Sync HashiCorp Vault KV v2 secrets to GitHub Actions secrets and environment variables.

## Prerequisites

- A running HashiCorp Vault instance with KV v2 secrets
- A GitHub repository with Actions enabled
- Vault token with read access to the source path
- GitHub personal access token with `repo` and `admin:org` scopes

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.vault_address` | Vault server URL |
| `source.vault_path` | Path within the KV v2 mount |
| `source.token` | Vault token |
| `source.vault_mount` | (Optional) KV v2 mount point, defaults to `secret` |
| `destination.repo` | Target GitHub repository in `owner/repo` format |
| `destination.token` | GitHub personal access token |
| `destination.url` | (Optional) GitHub Enterprise API URL |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_VAULT_TOKEN=hvs.xxx \
SECRET_SHIFT_SRC_VAULT_ADDRESS=https://vault.example.com \
SECRET_SHIFT_SRC_VAULT_PATH=myapp/config \
SECRET_SHIFT_DST_GITHUB_TOKEN=ghp_xxx \
SECRET_SHIFT_DST_GITHUB_REPO=owner/repo \
secret-shift sync -c config.json
```
