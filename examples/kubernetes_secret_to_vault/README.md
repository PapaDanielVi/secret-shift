# Kubernetes Secret to Vault

Sync a Kubernetes Secret to a HashiCorp Vault KV v2 secret.

## Prerequisites

- A running Kubernetes cluster with the source Secret
- A running HashiCorp Vault instance
- `kubeconfig` configured with read access to the source namespace
- Vault token with write access to the target path

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.kube_namespace` | Kubernetes namespace containing the Secret |
| `source.kube_secret_name` | Name of the source Kubernetes Secret |
| `source.kube_config` | (Optional) Path to kubeconfig file |
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
SECRET_SHIFT_SRC_KUBERNETES_KUBE_NAMESPACE=default \
SECRET_SHIFT_SRC_KUBERNETES_KUBE_SECRET_NAME=my-app-secrets \
SECRET_SHIFT_DST_VAULT_TOKEN=hvs.xxx \
SECRET_SHIFT_DST_VAULT_ADDRESS=https://vault.example.com \
SECRET_SHIFT_DST_VAULT_PATH=myapp/config \
secret-shift sync -c config.json
```
