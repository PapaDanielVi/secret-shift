# GitHub to Vault

Migrate GitHub Actions secrets and environment variables to a HashiCorp Vault KV v2 secret.

## Prerequisites

- A GitHub repository with Actions secrets/variables
- A running HashiCorp Vault instance
- Vault token with read/write access to the target KV v2 mount

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.repo` | Source GitHub repository in `owner/repo` format |
| `source.token` | GitHub personal access token |
| `source.url` | (Optional) GitHub Enterprise API URL |
| `destination.vault_address` | Vault server URL |
| `destination.vault_path` | Path within the KV v2 mount |
| `destination.token` | Vault token |
| `destination.vault_mount` | (Optional) KV v2 mount point, defaults to `secret` |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_GITHUB_TOKEN=ghp_xxx \
SECRET_SHIFT_SRC_GITHUB_REPO=owner/repo \
SECRET_SHIFT_DST_VAULT_TOKEN=hvs.xxx \
SECRET_SHIFT_DST_VAULT_ADDRESS=https://vault.example.com \
SECRET_SHIFT_DST_VAULT_PATH=myapp/config \
secret-shift sync -c config.json
```
