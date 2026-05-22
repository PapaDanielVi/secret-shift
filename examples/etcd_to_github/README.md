# `etcd` to `github` Example

Reads key-value pairs from an etcd cluster and syncs them as GitHub Actions secrets in a repository.

## Prerequisites

- A running etcd cluster accessible from the machine running secret-shift
- Optional: etcd authentication credentials if the cluster requires them
- A GitHub personal access token (PAT) with `repo` or `admin:org` scope
- The GitHub repository must exist

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
| `destination.type` | Must be `"github"` |
| `destination.repo` | Target repository in `"owner/repo"` format |
| `destination.token_redacted` | Your GitHub personal access token |
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
| `SECRET_SHIFT_DST_REPO` | `destination.repo` |
| `SECRET_SHIFT_DST_TOKEN` | `destination.token_redacted` |
| `SECRET_SHIFT_DST_CONFLICT_STRATEGY` | `destination.conflict_strategy` |
| `SECRET_SHIFT_DRY_RUN` | `dry_run` |
