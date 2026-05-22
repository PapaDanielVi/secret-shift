# GitHub to GitHub

Copy GitHub Actions secrets and environment variables from one repository to another.

## Prerequisites

- Source and destination GitHub repositories
- A GitHub personal access token with `repo` and `admin:org` scopes

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.repo` | Source repository in `owner/repo` format |
| `source.token` | GitHub personal access token |
| `source.url` | (Optional) GitHub Enterprise API URL |
| `destination.repo` | Destination repository in `owner/repo` format |
| `destination.token` | GitHub personal access token (can be same as source) |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_GITHUB_TOKEN=ghp_xxx \
SECRET_SHIFT_SRC_GITHUB_REPO=owner/source-repo \
SECRET_SHIFT_DST_GITHUB_TOKEN=ghp_xxx \
SECRET_SHIFT_DST_GITHUB_REPO=owner/dest-repo \
secret-shift sync -c config.json
```
