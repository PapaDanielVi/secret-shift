# Vault to Kubernetes ConfigMap

Sync HashiCorp Vault KV v2 secrets to a Kubernetes ConfigMap.

## Prerequisites

- A running HashiCorp Vault instance with KV v2 secrets
- A running Kubernetes cluster accessible via kubectl
- Vault token with read access to the source path
- `kubeconfig` configured with access to the target namespace

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.vault_address` | Vault server URL |
| `source.vault_path` | Path within the KV v2 mount |
| `source.token` | Vault token |
| `source.vault_mount` | (Optional) KV v2 mount point, defaults to `secret` |
| `destination.kube_namespace` | Target Kubernetes namespace |
| `destination.kube_secret_name` | Name of the Kubernetes ConfigMap to create/update |
| `destination.kube_config` | (Optional) Path to kubeconfig file |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_VAULT_TOKEN=hvs.xxx \
SECRET_SHIFT_SRC_VAULT_ADDRESS=https://vault.example.com \
SECRET_SHIFT_SRC_VAULT_PATH=myapp/config \
SECRET_SHIFT_DST_KUBERNETES_KUBE_NAMESPACE=default \
SECRET_SHIFT_DST_KUBERNETES_KUBE_SECRET_NAME=my-app-config \
secret-shift sync -c config.json
```
