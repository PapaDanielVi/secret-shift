# `kubernetes_configmap` to `kubernetes_configmap` Example

Reads key-value pairs from a Kubernetes ConfigMap and writes them into another Kubernetes ConfigMap, optionally in a different namespace.

## Prerequisites

- A running Kubernetes cluster accessible from your machine
- `kubectl` configured with access to both source and destination namespaces
- The source ConfigMap must exist in the specified source namespace
- Sufficient RBAC permissions to create/update ConfigMaps in the destination namespace

## Configuration

Edit `config.json` with your values:

| Field | Description |
|---|---|
| `source.type` | Must be `"kubernetes"` |
| `source.kube_namespace` | Kubernetes namespace where the source ConfigMap lives (e.g. `"default"`) |
| `source.kube_secret_name` | Name of the source Kubernetes ConfigMap to read |
| `process` | Optional filtering/transformation rules (leave empty to copy all keys) |
| `destination.type` | Must be `"kubernetes"` |
| `destination.kube_namespace` | Target Kubernetes namespace for the ConfigMap (e.g. `"production"`) |
| `destination.kube_secret_name` | Name of the target Kubernetes ConfigMap to create or update |
| `destination.conflict_strategy` | How to handle existing keys: `"replace"` or `"skip"` |
| `dry_run` | Set to `true` to preview changes without writing |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

All fields can be set via environment variables using the `SECRET_SHIFT_SRC_` and `SECRET_SHIFT_DST_` prefixes:

| Variable | Maps to |
|---|---|
| `SECRET_SHIFT_SRC_TYPE` | `source.type` |
| `SECRET_SHIFT_SRC_KUBE_NAMESPACE` | `source.kube_namespace` |
| `SECRET_SHIFT_SRC_KUBE_SECRET_NAME` | `source.kube_secret_name` |
| `SECRET_SHIFT_DST_TYPE` | `destination.type` |
| `SECRET_SHIFT_DST_KUBE_NAMESPACE` | `destination.kube_namespace` |
| `SECRET_SHIFT_DST_KUBE_SECRET_NAME` | `destination.kube_secret_name` |
| `SECRET_SHIFT_DST_CONFLICT_STRATEGY` | `destination.conflict_strategy` |
| `SECRET_SHIFT_DRY_RUN` | `dry_run` |
