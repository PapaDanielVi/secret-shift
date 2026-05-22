# `kubernetes_secret` to `file` Example

Reads secrets from a Kubernetes Secret resource and writes them to a local JSON or YAML file.

## Prerequisites

- A running Kubernetes cluster accessible from your machine
- `kubectl` configured with access to the target namespace
- The source Secret must exist in the specified namespace

## Configuration

Edit `config.json` with your values:

| Field | Description |
|---|---|
| `source.type` | Must be `"kubernetes"` |
| `source.kube_namespace` | Kubernetes namespace where the Secret lives (e.g. `"default"`) |
| `source.kube_secret_name` | Name of the Kubernetes Secret to read |
| `process` | Optional filtering/transformation rules (leave empty to copy all keys) |
| `destination.type` | Must be `"file"` |
| `destination.path` | Output file path (e.g. `"./output/k8s-secrets.json"`) |
| `destination.format` | Output format: `"json"` or `"yaml"` |
| `destination.encrypt` | Set to `true` to encrypt the output file |
| `destination.encrypt_key` | Encryption key (required if `encrypt` is `true`) |
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
| `SECRET_SHIFT_DST_PATH` | `destination.path` |
| `SECRET_SHIFT_DST_FORMAT` | `destination.format` |
| `SECRET_SHIFT_DST_ENCRYPT` | `destination.encrypt` |
| `SECRET_SHIFT_DST_ENCRYPT_KEY` | `destination.encrypt_key` |
| `SECRET_SHIFT_DRY_RUN` | `dry_run` |
