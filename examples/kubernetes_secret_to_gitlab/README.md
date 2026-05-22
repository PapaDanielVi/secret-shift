# Kubernetes Secret to GitLab

Sync a Kubernetes Secret to GitLab project CI/CD variables.

## Prerequisites

- A running Kubernetes cluster with the source Secret
- A GitLab project
- `kubeconfig` configured with read access to the source namespace
- GitLab personal access token with `api` scope

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.kube_namespace` | Kubernetes namespace containing the Secret |
| `source.kube_secret_name` | Name of the source Kubernetes Secret |
| `source.kube_config` | (Optional) Path to kubeconfig file |
| `destination.project_id` | Target GitLab project ID |
| `destination.token` | GitLab personal access token |
| `destination.url` | (Optional) Self-hosted GitLab URL |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_KUBERNETES_KUBE_NAMESPACE=default \
SECRET_SHIFT_SRC_KUBERNETES_KUBE_SECRET_NAME=my-app-secrets \
SECRET_SHIFT_DST_GITLAB_TOKEN=glpat-xxx \
SECRET_SHIFT_DST_GITLAB_PROJECT_ID=123 \
secret-shift sync -c config.json
```
