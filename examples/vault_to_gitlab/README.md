# Vault to GitLab

Sync HashiCorp Vault KV v2 secrets to GitLab project CI/CD variables.

## Prerequisites

- A running HashiCorp Vault instance with KV v2 secrets
- A GitLab project
- Vault token with read access to the source path
- GitLab personal access token with `api` scope

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.vault_address` | Vault server URL |
| `source.vault_path` | Path within the KV v2 mount |
| `source.token` | Vault token |
| `source.vault_mount` | (Optional) KV v2 mount point, defaults to `secret` |
| `destination.project_id` | Target GitLab project ID |
| `destination.token` | GitLab personal access token |
| `destination.url` | (Optional) Self-hosted GitLab URL |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_VAULT_TOKEN=hvs.xxx \
SECRET_SHIFT_SRC_VAULT_ADDRESS=https://vault.example.com \
SECRET_SHIFT_SRC_VAULT_PATH=myapp/config \
SECRET_SHIFT_DST_GITLAB_TOKEN=glpat-xxx \
SECRET_SHIFT_DST_GITLAB_PROJECT_ID=123 \
secret-shift sync -c config.json
```
