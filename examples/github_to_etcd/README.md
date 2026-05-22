# GitHub to etcd

Sync GitHub Actions secrets and environment variables to an etcd key-value store.

## Prerequisites

- A GitHub repository with Actions secrets/variables
- A running etcd instance
- Network access to the etcd endpoints

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.repo` | Source GitHub repository in `owner/repo` format |
| `source.token` | GitHub personal access token |
| `source.url` | (Optional) GitHub Enterprise API URL |
| `destination.etcd_endpoints` | Array of etcd endpoint URLs |
| `destination.etcd_prefix` | Key prefix for stored values |
| `destination.etcd_username` | (Optional) etcd username |
| `destination.etcd_password` | (Optional) etcd password |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_GITHUB_TOKEN=ghp_xxx \
SECRET_SHIFT_SRC_GITHUB_REPO=owner/repo \
SECRET_SHIFT_DST_ETCD_ENDPOINTS='["http://etcd:2379"]' \
SECRET_SHIFT_DST_ETCD_PREFIX=/myapp/ \
secret-shift sync -c config.json
```
