# GitHub to File

Export GitHub Actions secrets and environment variables to a local JSON or YAML file.

## Prerequisites

- A GitHub repository with Actions secrets/variables
- Write access to the output directory

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.repo` | Source GitHub repository in `owner/repo` format |
| `source.token` | GitHub personal access token |
| `source.url` | (Optional) GitHub Enterprise API URL |
| `destination.path` | Output file path |
| `destination.format` | Output format: `json` or `yaml` |
| `destination.encrypt` | Set to `true` for AES-256-GCM encryption |
| `destination.encrypt_key` | Encryption passphrase (required if encrypt is true) |

## Run

```bash
secret-shift sync -c config.json
```

## Environment Variables

```bash
SECRET_SHIFT_SRC_GITHUB_TOKEN=ghp_xxx \
SECRET_SHIFT_SRC_GITHUB_REPO=owner/repo \
SECRET_SHIFT_DST_FILE_PATH=./output/secrets.json \
SECRET_SHIFT_DST_FILE_FORMAT=json \
secret-shift sync -c config.json
```
