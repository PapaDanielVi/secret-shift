# Kubernetes Secret to Kubernetes ConfigMap

Convert a Kubernetes Secret to a Kubernetes ConfigMap.

## Prerequisites

- A running Kubernetes cluster with the source Secret
- `kubeconfig` configured with access to the target namespace

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.kube_namespace` | Kubernetes namespace containing the Secret |
| `source.kube_secret_name` | Name of the source Kubernetes Secret |
| `source.kube_config` | (Optional) Path to kubeconfig file |
| `destination.kube_namespace` | Target Kubernetes namespace |
| `destination.kube_secret_name` | Name of the Kubernetes ConfigMap to create/update |
| `destination.kube_config` | (Optional) Path to kubeconfig file |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_KUBERNETES_KUBE_NAMESPACE=default \
SECRET_SHIFT_SRC_KUBERNETES_KUBE_SECRET_NAME=my-app-secrets \
SECRET_SHIFT_DST_KUBERNETES_KUBE_NAMESPACE=default \
SECRET_SHIFT_DST_KUBERNETES_KUBE_SECRET_NAME=my-app-config \
secret-shift sync -c config.json
```
