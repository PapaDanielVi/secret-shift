# Vault to File

Export HashiCorp Vault KV v2 secrets to a local JSON or YAML file.

## Prerequisites

- A running HashiCorp Vault instance with KV v2 secrets
- Vault token with read access to the source path
- Write access to the output directory

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.vault_address` | Vault server URL |
| `source.vault_path` | Path within the KV v2 mount |
| `source.token` | Vault token |
| `source.vault_mount` | (Optional) KV v2 mount point, defaults to `secret` |
| `destination.path` | Output file path |
| `destination.format` | Output format: `json` or `yaml` |
| `destination.encrypt` | Set to `true` for AES-256-GCM encryption |
| `destination.encrypt_key` | Encryption passphrase (required if encrypt is true) |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_VAULT_TOKEN=hvs.xxx \
SECRET_SHIFT_SRC_VAULT_ADDRESS=https://vault.example.com \
SECRET_SHIFT_SRC_VAULT_PATH=myapp/config \
SECRET_SHIFT_DST_FILE_PATH=./output/secrets.json \
SECRET_SHIFT_DST_FILE_FORMAT=json \
secret-shift sync -c config.json
```
