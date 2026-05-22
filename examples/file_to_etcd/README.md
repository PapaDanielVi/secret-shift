# `file` to `etcd` Example

Reads key-value pairs from a local JSON or YAML file and writes them to an etcd cluster.

## Prerequisites

- A source file in JSON or YAML format containing key-value pairs
- A running etcd cluster accessible from the machine running secret-shift
- Optional: etcd authentication credentials if the cluster requires them

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
| `SECRET_SHIFT_SRC_PATH` | `source.path` |
| `SECRET_SHIFT_SRC_FORMAT` | `source.format` |
| `SECRET_SHIFT_SRC_ENCRYPT` | `source.encrypt` |
| `SECRET_SHIFT_SRC_ENCRYPT_KEY` | `source.encrypt_key` |
| `SECRET_SHIFT_DST_TYPE` | `destination.type` |
| `SECRET_SHIFT_DST_ETCD_ENDPOINTS` | `destination.etcd_endpoints` |
| `SECRET_SHIFT_DST_ETCD_PREFIX` | `destination.etcd_prefix` |
| `SECRET_SHIFT_DST_ETCD_USERNAME` | `destination.etcd_username` |
| `SECRET_SHIFT_DST_ETCD_PASSWORD` | `destination.etcd_password` |
| `SECRET_SHIFT_DST_CONFLICT_STRATEGY` | `destination.conflict_strategy` |
| `SECRET_SHIFT_DRY_RUN` | `dry_run` |
