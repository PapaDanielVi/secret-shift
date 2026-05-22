# Kubernetes Secret to Kubernetes Secret

Copy a Kubernetes Secret to another namespace or cluster.

## Prerequisites

- Source and destination Kubernetes clusters (can be the same)
- `kubeconfig` configured with appropriate access to both namespaces

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.kube_namespace` | Source Kubernetes namespace |
| `source.kube_secret_name` | Name of the source Kubernetes Secret |
| `source.kube_config` | (Optional) Path to kubeconfig for source cluster |
| `destination.kube_namespace` | Destination Kubernetes namespace |
| `destination.kube_secret_name` | Name of the destination Kubernetes Secret |
| `destination.kube_config` | (Optional) Path to kubeconfig for destination cluster |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_KUBERNETES_KUBE_NAMESPACE=default \
SECRET_SHIFT_SRC_KUBERNETES_KUBE_SECRET_NAME=my-app-secrets \
SECRET_SHIFT_DST_KUBERNETES_KUBE_NAMESPACE=production \
SECRET_SHIFT_DST_KUBERNETES_KUBE_SECRET_NAME=my-app-secrets \
secret-shift sync -c config.json
```
