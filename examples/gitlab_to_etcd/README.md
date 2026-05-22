# GitLab to etcd

Sync GitLab project CI/CD variables to an etcd key-value store.

## Prerequisites

- A GitLab project with CI/CD variables
- A running etcd instance
- Network access to the etcd endpoints

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.project_id` | Your GitLab project ID |
| `source.token` | GitLab personal access token |
| `source.url` | (Optional) Self-hosted GitLab URL |
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
SECRET_SHIFT_SRC_GITLAB_TOKEN=glpat-xxx \
SECRET_SHIFT_SRC_GITLAB_PROJECT_ID=123 \
SECRET_SHIFT_DST_ETCD_ENDPOINTS='["http://etcd:2379"]' \
SECRET_SHIFT_DST_ETCD_PREFIX=/myapp/ \
secret-shift sync -c config.json
```
