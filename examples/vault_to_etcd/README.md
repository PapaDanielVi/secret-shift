# Vault to etcd

Sync HashiCorp Vault KV v2 secrets to an etcd key-value store.

## Prerequisites

- A running HashiCorp Vault instance with KV v2 secrets
- A running etcd instance
- Vault token with read access to the source path
- Network access to the etcd endpoints

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.vault_address` | Vault server URL |
| `source.vault_path` | Path within the KV v2 mount |
| `source.token` | Vault token |
| `source.vault_mount` | (Optional) KV v2 mount point, defaults to `secret` |
| `destination.etcd_endpoints` | Array of etcd endpoint URLs |
| `destination.etcd_prefix` | Key prefix for stored values |
| `destination.etcd_username` | (Optional) etcd username |
| `destination.etcd_password` | (Optional) etcd password |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_VAULT_TOKEN=hvs.xxx \
SECRET_SHIFT_SRC_VAULT_ADDRESS=https://vault.example.com \
SECRET_SHIFT_SRC_VAULT_PATH=myapp/config \
SECRET_SHIFT_DST_ETCD_ENDPOINTS='["http://etcd:2379"]' \
SECRET_SHIFT_DST_ETCD_PREFIX=/myapp/ \
secret-shift sync -c config.json
```
