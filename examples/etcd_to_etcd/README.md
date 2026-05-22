# `etcd` to `etcd` Example

Reads key-value pairs from one etcd cluster and writes them to another etcd cluster (or a different prefix in the same cluster). Useful for cross-cluster replication.

## Prerequisites

- Source and destination etcd clusters accessible from the machine running secret-shift
- Optional: etcd authentication credentials for either or both clusters

## Configuration

Edit `config.json` with your values:

| Field | Description |
|---|---|
| `source.type` | Must be `"etcd"` |
| `source.etcd_endpoints` | Array of source etcd endpoints (e.g. `["http://localhost:2379"]`) |
| `source.etcd_prefix` | Key prefix to read from (e.g. `"/my-app/"`) |
| `source.etcd_username` | Source etcd username (optional) |
| `source.etcd_password` | Source etcd password (optional) |
| `process` | Optional filtering/transformation rules (leave empty to copy all keys) |
| `destination.type` | Must be `"etcd"` |
| `destination.etcd_endpoints` | Array of destination etcd endpoints (e.g. `["http://remote:2379"]`) |
| `destination.etcd_prefix` | Key prefix for stored values (e.g. `"/my-app/"`) |
| `destination.etcd_username` | Destination etcd username (optional) |
| `destination.etcd_password` | Destination etcd password (optional) |
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
| `SECRET_SHIFT_SRC_ETCD_ENDPOINTS` | `source.etcd_endpoints` |
| `SECRET_SHIFT_SRC_ETCD_PREFIX` | `source.etcd_prefix` |
| `SECRET_SHIFT_SRC_ETCD_USERNAME` | `source.etcd_username` |
| `SECRET_SHIFT_SRC_ETCD_PASSWORD` | `source.etcd_password` |
| `SECRET_SHIFT_DST_TYPE` | `destination.type` |
| `SECRET_SHIFT_DST_ETCD_ENDPOINTS` | `destination.etcd_endpoints` |
| `SECRET_SHIFT_DST_ETCD_PREFIX` | `destination.etcd_prefix` |
| `SECRET_SHIFT_DST_ETCD_USERNAME` | `destination.etcd_username` |
| `SECRET_SHIFT_DST_ETCD_PASSWORD` | `destination.etcd_password` |
| `SECRET_SHIFT_DST_CONFLICT_STRATEGY` | `destination.conflict_strategy` |
| `SECRET_SHIFT_DRY_RUN` | `dry_run` |
