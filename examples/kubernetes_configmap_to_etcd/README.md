# `kubernetes_configmap` to `etcd` Example

Reads key-value pairs from a Kubernetes ConfigMap and writes them to an etcd cluster.

## Prerequisites

- A running Kubernetes cluster accessible from your machine
- `kubectl` configured with access to the target namespace
- The source ConfigMap must exist in the specified namespace
- A running etcd cluster accessible from the machine running secret-shift
- Optional: etcd authentication credentials if the cluster requires them

## Configuration

Edit `config.json` with your values:

| Field | Description |
|---|---|
| `source.type` | Must be `"kubernetes"` |
| `source.kube_namespace` | Kubernetes namespace where the ConfigMap lives (e.g. `"default"`) |
| `source.kube_secret_name` | Name of the Kubernetes ConfigMap to read (used as the configmap name) |
| `process` | Optional filtering/transformation rules (leave empty to copy all keys) |
| `destination.type` | Must be `"etcd"` |
| `destination.etcd_endpoints` | Array of etcd endpoints (e.g. `["http://localhost:2379"]`) |
| `destination.etcd_prefix` | Key prefix for stored values (e.g. `"/my-app/"`) |
| `destination.etcd_username` | etcd username (optional, for authenticated clusters) |
| `destination.etcd_password` | etcd password (optional, for authenticated clusters) |
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
| `SECRET_SHIFT_DST_ETCD_ENDPOINTS` | `destination.etcd_endpoints` |
| `SECRET_SHIFT_DST_ETCD_PREFIX` | `destination.etcd_prefix` |
| `SECRET_SHIFT_DST_ETCD_USERNAME` | `destination.etcd_username` |
| `SECRET_SHIFT_DST_ETCD_PASSWORD` | `destination.etcd_password` |
| `SECRET_SHIFT_DST_CONFLICT_STRATEGY` | `destination.conflict_strategy` |
| `SECRET_SHIFT_DRY_RUN` | `dry_run` |
