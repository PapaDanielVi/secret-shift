# GitLab to GitLab

Copy GitLab project CI/CD variables from one project to another.

## Prerequisites

- Source and destination GitLab projects
- A GitLab personal access token with `api` scope for both projects

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.project_id` | Source GitLab project ID |
| `source.token` | GitLab personal access token |
| `source.url` | (Optional) Self-hosted GitLab URL |
| `destination.project_id` | Destination GitLab project ID |
| `destination.token` | GitLab personal access token (can be same as source) |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_GITLAB_TOKEN=glpat-xxx \
SECRET_SHIFT_SRC_GITLAB_PROJECT_ID=123 \
SECRET_SHIFT_DST_GITLAB_TOKEN=glpat-xxx \
SECRET_SHIFT_DST_GITLAB_PROJECT_ID=456 \
secret-shift sync -c config.json
```
