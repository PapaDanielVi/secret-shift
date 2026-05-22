# `etcd` to `kubernetes_secret` Example

Reads key-value pairs from an etcd cluster and writes them into a Kubernetes Secret.

## Prerequisites

- A running etcd cluster accessible from the machine running secret-shift
- Optional: etcd authentication credentials if the cluster requires them
- A running Kubernetes cluster accessible from your machine
- `kubectl` configured with access to the target namespace
- Sufficient RBAC permissions to create/update Secrets in the destination namespace

## Configuration

Edit `config.json` with your values:

| Field | Description |
|---|---|
| `source.type` | Must be `"etcd"` |
| `source.etcd_endpoints` | Array of etcd endpoints (e.g. `["http://localhost:2379"]`) |
| `source.etcd_prefix` | Key prefix to read from (e.g. `"/my-app/"`) |
| `source.etcd_username` | etcd username (optional, for authenticated clusters) |
| `source.etcd_password` | etcd password (optional, for authenticated clusters) |
| `process` | Optional filtering/transformation rules (leave empty to copy all keys) |
| `destination.type` | Must be `"kubernetes"` |
| `destination.kube_namespace` | Target Kubernetes namespace for the Secret (e.g. `"default"`) |
| `destination.kube_secret_name` | Name of the Kubernetes Secret to create or update |
| `destination.conflict_strategy` | How to handle existing secrets: `"replace"` or `"skip"` |
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
| `SECRET_SHIFT_SRC_ETCD_ENDPOINTS` | `source.etcd_endpoints` |
| `SECRET_SHIFT_SRC_ETCD_PREFIX` | `source.etcd_prefix` |
| `SECRET_SHIFT_SRC_ETCD_USERNAME` | `source.etcd_username` |
| `SECRET_SHIFT_SRC_ETCD_PASSWORD` | `source.etcd_password` |
| `SECRET_SHIFT_DST_TYPE` | `destination.type` |
| `SECRET_SHIFT_DST_KUBE_NAMESPACE` | `destination.kube_namespace` |
| `SECRET_SHIFT_DST_KUBE_SECRET_NAME` | `destination.kube_secret_name` |
| `SECRET_SHIFT_DST_CONFLICT_STRATEGY` | `destination.conflict_strategy` |
| `SECRET_SHIFT_DRY_RUN` | `dry_run` |
