# `file` to `kubernetes_configmap` Example

Reads key-value pairs from a local JSON or YAML file and writes them into a Kubernetes ConfigMap.

## Prerequisites

- A source file in JSON or YAML format containing key-value pairs
- A running Kubernetes cluster accessible from your machine
- `kubectl` configured with access to the target namespace
- Sufficient RBAC permissions to create/update ConfigMaps in the destination namespace

## Configuration

Edit `config.json` with your values:

| Field | Description |
|---|---|
| `source.type` | Must be `"file"` |
| `source.path` | Path to the source file (e.g. `"./secrets.json"`) |
| `source.format` | File format: `"json"` or `"yaml"` |
| `source.encrypt` | Set to `true` if the source file is encrypted |
| `source.encrypt_key` | Decryption key (required if `encrypt` is `true`) |
| `process` | Optional filtering/transformation rules (leave empty to copy all keys) |
| `destination.type` | Must be `"kubernetes"` |
| `destination.kube_namespace` | Target Kubernetes namespace for the ConfigMap (e.g. `"default"`) |
| `destination.kube_secret_name` | Name of the Kubernetes ConfigMap to create or update |
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
| `SECRET_SHIFT_SRC_PATH` | `source.path` |
| `SECRET_SHIFT_SRC_FORMAT` | `source.format` |
| `SECRET_SHIFT_SRC_ENCRYPT` | `source.encrypt` |
| `SECRET_SHIFT_SRC_ENCRYPT_KEY` | `source.encrypt_key` |
| `SECRET_SHIFT_DST_TYPE` | `destination.type` |
| `SECRET_SHIFT_DST_KUBE_NAMESPACE` | `destination.kube_namespace` |
| `SECRET_SHIFT_DST_KUBE_SECRET_NAME` | `destination.kube_secret_name` |
| `SECRET_SHIFT_DST_CONFLICT_STRATEGY` | `destination.conflict_strategy` |
| `SECRET_SHIFT_DRY_RUN` | `dry_run` |
