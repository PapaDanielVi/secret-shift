# `file` to `gitlab` Example

Reads key-value pairs from a local JSON or YAML file and syncs them as GitLab CI/CD variables in a project.

## Prerequisites

- A source file in JSON or YAML format containing key-value pairs
- A GitLab personal access token with `api` scope
- The GitLab project must exist

## Configuration

Edit `config.json` with your values:

| Field | Description |
|---|---|
| `source.type` | Must be `"file"` |
| `source.path` | Path to the source file (e.g. `"./secrets.json"`) |
| `source.format` | File format: `"json"` or `"yaml"` |
| `source.encrypt` | Set to `true` if the source file is encrypted |
| `source.encrypt_key` | Decryption key (required if `encrypt` is `true`) |
| `process` | Optional filtering/transformation rules (leave empty to copy all keys) |
| `destination.type` | Must be `"gitlab"` |
| `destination.project_id` | GitLab project ID (numeric) |
| `destination.token_redacted` | Your GitLab personal access token |
| `destination.url` | GitLab instance URL (optional, defaults to `https://gitlab.com`) |
| `destination.conflict_strategy` | How to handle existing variables: `"replace"` or `"skip"` |
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
| `SECRET_SHIFT_SRC_PATH` | `source.path` |
| `SECRET_SHIFT_SRC_FORMAT` | `source.format` |
| `SECRET_SHIFT_SRC_ENCRYPT` | `source.encrypt` |
| `SECRET_SHIFT_SRC_ENCRYPT_KEY` | `source.encrypt_key` |
| `SECRET_SHIFT_DST_TYPE` | `destination.type` |
| `SECRET_SHIFT_DST_PROJECT_ID` | `destination.project_id` |
| `SECRET_SHIFT_DST_TOKEN` | `destination.token_redacted` |
| `SECRET_SHIFT_DST_URL` | `destination.url` |
| `SECRET_SHIFT_DST_CONFLICT_STRATEGY` | `destination.conflict_strategy` |
| `SECRET_SHIFT_DRY_RUN` | `dry_run` |
