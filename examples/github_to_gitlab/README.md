# GitHub to GitLab

Sync GitHub Actions secrets and environment variables to GitLab project CI/CD variables.

## Prerequisites

- A GitHub repository with Actions secrets/variables
- A GitLab project
- A GitHub personal access token with `repo` and `admin:org` scopes
- A GitLab personal access token with `api` scope

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.repo` | Source GitHub repository in `owner/repo` format |
| `source.token` | GitHub personal access token |
| `source.url` | (Optional) GitHub Enterprise API URL |
| `destination.project_id` | Target GitLab project ID |
| `destination.token` | GitLab personal access token |
| `destination.url` | (Optional) Self-hosted GitLab URL |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_GITHUB_TOKEN=ghp_xxx \
SECRET_SHIFT_SRC_GITHUB_REPO=owner/repo \
SECRET_SHIFT_DST_GITLAB_TOKEN=glpat-xxx \
SECRET_SHIFT_DST_GITLAB_PROJECT_ID=123 \
secret-shift sync -c config.json
```
