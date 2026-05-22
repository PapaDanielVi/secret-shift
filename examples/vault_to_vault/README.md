# Vault to Vault

Copy secrets between paths in HashiCorp Vault (same or different instances).

## Prerequisites

- Source and destination Vault instances (can be the same)
- Vault tokens with read access to source and write access to destination

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.vault_address` | Source Vault server URL |
| `source.vault_path` | Source path within the KV v2 mount |
| `source.token` | Vault token for source |
| `source.vault_mount` | (Optional) KV v2 mount point, defaults to `secret` |
| `destination.vault_address` | Destination Vault server URL |
| `destination.vault_path` | Destination path within the KV v2 mount |
| `destination.token` | Vault token for destination |
| `destination.vault_mount` | (Optional) KV v2 mount point, defaults to `secret` |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_VAULT_TOKEN=hvs.xxx \
SECRET_SHIFT_SRC_VAULT_ADDRESS=https://vault.example.com \
SECRET_SHIFT_SRC_VAULT_PATH=myapp/source-config \
SECRET_SHIFT_DST_VAULT_TOKEN=hvs.yyy \
SECRET_SHIFT_DST_VAULT_ADDRESS=https://vault.example.com \
SECRET_SHIFT_DST_VAULT_PATH=myapp/dest-config \
secret-shift sync -c config.json
```
