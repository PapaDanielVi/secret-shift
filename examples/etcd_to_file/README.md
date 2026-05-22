# `etcd` to `file` Example

Reads key-value pairs from an etcd cluster and writes them to a local JSON or YAML file.

## Prerequisites

- A running etcd cluster accessible from the machine running secret-shift
- Optional: etcd authentication credentials if the cluster requires them

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
| `destination.type` | Must be `"file"` |
| `destination.path` | Output file path (e.g. `"./output/etcd-secrets.json"`) |
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
| `SECRET_SHIFT_SRC_ETCD_ENDPOINTS` | `source.etcd_endpoints` |
| `SECRET_SHIFT_SRC_ETCD_PREFIX` | `source.etcd_prefix` |
| `SECRET_SHIFT_SRC_ETCD_USERNAME` | `source.etcd_username` |
| `SECRET_SHIFT_SRC_ETCD_PASSWORD` | `source.etcd_password` |
| `SECRET_SHIFT_DST_TYPE` | `destination.type` |
| `SECRET_SHIFT_DST_PATH` | `destination.path` |
| `SECRET_SHIFT_DST_FORMAT` | `destination.format` |
| `SECRET_SHIFT_DST_ENCRYPT` | `destination.encrypt` |
| `SECRET_SHIFT_DST_ENCRYPT_KEY` | `destination.encrypt_key` |
| `SECRET_SHIFT_DRY_RUN` | `dry_run` |
