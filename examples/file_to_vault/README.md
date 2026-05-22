# `file` to `vault` Example

Reads key-value pairs from a local JSON or YAML file and writes them to a HashiCorp Vault secret path.

## Prerequisites

- A source file in JSON or YAML format containing key-value pairs
- A running Vault server with the KV v2 secrets engine enabled
- A Vault token with write access to the target path

## Configuration

Edit `config.json` with your values:

| Field | Description |
|---|---|
| `source.type` | Must be `"file"` |
| `source.path` | Path to the source file (e.g. `"./secrets.json"`) |
| `source.format` | File format: `"json"` or `"yaml"` |
| `source.encrypt` | Set to `true` if the source file is encrypted |
| `source.encrypt_key` | Decryption key (required if `encrypt` is `true`) |
| `process` | Optional filtering/transformation rules (leave empty to copy all keys) |
| `destination.type` | Must be `"vault"` |
| `destination.vault_address` | Vault server URL (e.g. `"https://vault.example.com"`) |
| `destination.vault_path` | Path to the secret in Vault (e.g. `"secret/data/my-app"`) |
| `destination.vault_mount` | Vault secrets engine mount point (optional, default `"secret"`) |
| `destination.token_redacted` | Your Vault token |
| `destination.conflict_strategy` | How to handle existing secrets: `"replace"` or `"skip"` |
| `dry_run` | Set to `true` to preview changes without writing |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

All fields can be set via environment variables using the `SECRET_SHIFT_SRC_` and `SECRET_SHIFT_DST_` prefixes:

| Variable | Maps to |
|---|---|
| `SECRET_SHIFT_SRC_TYPE` | `source.type` |
| `SECRET_SHIFT_SRC_PATH` | `source.path` |
| `SECRET_SHIFT_SRC_FORMAT` | `source.format` |
| `SECRET_SHIFT_SRC_ENCRYPT` | `source.encrypt` |
| `SECRET_SHIFT_SRC_ENCRYPT_KEY` | `source.encrypt_key` |
| `SECRET_SHIFT_DST_TYPE` | `destination.type` |
| `SECRET_SHIFT_DST_VAULT_ADDRESS` | `destination.vault_address` |
| `SECRET_SHIFT_DST_VAULT_PATH` | `destination.vault_path` |
| `SECRET_SHIFT_DST_VAULT_MOUNT` | `destination.vault_mount` |
| `SECRET_SHIFT_DST_TOKEN` | `destination.token_redacted` |
| `SECRET_SHIFT_DST_CONFLICT_STRATEGY` | `destination.conflict_strategy` |
| `SECRET_SHIFT_DRY_RUN` | `dry_run` |
