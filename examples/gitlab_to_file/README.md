# GitLab to File

Export GitLab project CI/CD variables to a local JSON or YAML file.

## Prerequisites

- A GitLab project with CI/CD variables
- Write access to the output directory

## Configuration

Replace the placeholder values in `config.json`:

| Field | Description |
| ----- | ----------- |
| `source.project_id` | Your GitLab project ID |
| `source.token` | GitLab personal access token |
| `source.url` | (Optional) Self-hosted GitLab URL |
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
SECRET_SHIFT_SRC_GITLAB_TOKEN=glpat-xxx \
SECRET_SHIFT_SRC_GITLAB_PROJECT_ID=123 \
SECRET_SHIFT_DST_FILE_PATH=./output/secrets.json \
SECRET_SHIFT_DST_FILE_FORMAT=json \
secret-shift sync -c config.json
```
