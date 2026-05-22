# GitHub to Kubernetes ConfigMap

Sync GitHub Actions secrets and environment variables to a Kubernetes ConfigMap.

## Prerequisites

- A GitHub repository with Actions secrets/variables
- A running Kubernetes cluster accessible via kubectl
- `kubeconfig` configured with access to the target namespace

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.repo` | Source GitHub repository in `owner/repo` format |
| `source.token` | GitHub personal access token |
| `source.url` | (Optional) GitHub Enterprise API URL |
| `destination.kube_namespace` | Target Kubernetes namespace |
| `destination.kube_secret_name` | Name of the Kubernetes ConfigMap to create/update |
| `destination.kube_config` | (Optional) Path to kubeconfig file |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_GITHUB_TOKEN=ghp_xxx \
SECRET_SHIFT_SRC_GITHUB_REPO=owner/repo \
SECRET_SHIFT_DST_KUBERNETES_KUBE_NAMESPACE=default \
SECRET_SHIFT_DST_KUBERNETES_KUBE_SECRET_NAME=my-app-config \
secret-shift sync -c config.json
```
