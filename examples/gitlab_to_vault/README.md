# GitLab to Vault

Migrate GitLab project CI/CD variables to a HashiCorp Vault KV v2 secret.

## Prerequisites

- A GitLab project with CI/CD variables
- A running HashiCorp Vault instance
- Vault token with read/write access to the target KV v2 mount

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.project_id` | Your GitLab project ID |
| `source.token` | GitLab personal access token |
| `source.url` | (Optional) Self-hosted GitLab URL |
| `destination.vault_address` | Vault server URL |
| `destination.vault_path` | Path within the KV v2 mount (e.g. `myapp/config`) |
| `destination.token` | Vault token |
| `destination.vault_mount` | (Optional) KV v2 mount point, defaults to `secret` |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_GITLAB_TOKEN=glpat-xxx \
SECRET_SHIFT_SRC_GITLAB_PROJECT_ID=123 \
SECRET_SHIFT_DST_VAULT_TOKEN=hvs.xxx \
SECRET_SHIFT_DST_VAULT_ADDRESS=https://vault.example.com \
SECRET_SHIFT_DST_VAULT_PATH=myapp/config \
secret-shift sync -c config.json
```
