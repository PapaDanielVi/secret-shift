# Kubernetes Secret to GitHub

Sync a Kubernetes Secret to GitHub Actions secrets and environment variables.

## Prerequisites

- A running Kubernetes cluster with the source Secret
- A GitHub repository with Actions enabled
- `kubeconfig` configured with read access to the source namespace
- GitHub personal access token with `repo` and `admin:org` scopes

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.kube_namespace` | Kubernetes namespace containing the Secret |
| `source.kube_secret_name` | Name of the source Kubernetes Secret |
| `source.kube_config` | (Optional) Path to kubeconfig file |
| `destination.repo` | Target GitHub repository in `owner/repo` format |
| `destination.token` | GitHub personal access token |
| `destination.url` | (Optional) GitHub Enterprise API URL |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_KUBERNETES_KUBE_NAMESPACE=default \
SECRET_SHIFT_SRC_KUBERNETES_KUBE_SECRET_NAME=my-app-secrets \
SECRET_SHIFT_DST_GITHUB_TOKEN=ghp_xxx \
SECRET_SHIFT_DST_GITHUB_REPO=owner/repo \
secret-shift sync -c config.json
```
