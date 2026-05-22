# GitLab to Kubernetes Secret

Sync GitLab project CI/CD variables to a Kubernetes Secret.

## Prerequisites

- A GitLab project with CI/CD variables
- A running Kubernetes cluster accessible via kubectl
- `kubeconfig` configured with access to the target namespace

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.project_id` | Your GitLab project ID |
| `source.token` | GitLab personal access token |
| `source.url` | (Optional) Self-hosted GitLab URL |
| `destination.kube_namespace` | Target Kubernetes namespace |
| `destination.kube_secret_name` | Name of the Kubernetes Secret to create/update |
| `destination.kube_config` | (Optional) Path to kubeconfig file |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_GITLAB_TOKEN=glpat-xxx \
SECRET_SHIFT_SRC_GITLAB_PROJECT_ID=123 \
SECRET_SHIFT_DST_KUBERNETES_KUBE_NAMESPACE=default \
SECRET_SHIFT_DST_KUBERNETES_KUBE_SECRET_NAME=my-app-secrets \
secret-shift sync -c config.json
```
