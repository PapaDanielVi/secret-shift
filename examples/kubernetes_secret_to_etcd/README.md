# Kubernetes Secret to etcd

Sync a Kubernetes Secret to an etcd key-value store.

## Prerequisites

- A running Kubernetes cluster with the source Secret
- A running etcd instance
- `kubeconfig` configured with read access to the source namespace
- Network access to the etcd endpoints

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.kube_namespace` | Kubernetes namespace containing the Secret |
| `source.kube_secret_name` | Name of the source Kubernetes Secret |
| `source.kube_config` | (Optional) Path to kubeconfig file |
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
SECRET_SHIFT_SRC_KUBERNETES_KUBE_NAMESPACE=default \
SECRET_SHIFT_SRC_KUBERNETES_KUBE_SECRET_NAME=my-app-secrets \
SECRET_SHIFT_DST_ETCD_ENDPOINTS='["http://etcd:2379"]' \
SECRET_SHIFT_DST_ETCD_PREFIX=/myapp/ \
secret-shift sync -c config.json
```
