# GitLab to GitHub

Sync GitLab project CI/CD variables to GitHub Actions secrets and environment variables.

## Prerequisites

- A GitLab project with CI/CD variables
- A GitHub repository with Actions enabled
- A GitLab personal access token with `api` scope
- A GitHub personal access token with `repo` and `admin:org` scopes (for secrets)

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.project_id` | Your GitLab project ID (found on the project page) |
| `source.token` | GitLab personal access token |
| `source.url` | (Optional) Self-hosted GitLab URL |
| `destination.repo` | Target GitHub repository in `owner/repo` format |
| `destination.token` | GitHub personal access token |
| `destination.url` | (Optional) GitHub Enterprise API URL |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

All config values can be set via environment variables:

```bash
SECRET_SHIFT_SRC_GITLAB_TOKEN=glpat-xxx \
SECRET_SHIFT_SRC_GITLAB_PROJECT_ID=123 \
SECRET_SHIFT_DST_GITHUB_TOKEN=ghp_xxx \
SECRET_SHIFT_DST_GITHUB_REPO=owner/repo \
secret-shift sync -c config.json
```
