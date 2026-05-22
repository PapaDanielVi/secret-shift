# `file` to `file` Example

Reads key-value pairs from a local JSON or YAML file and writes them to another local file. Useful for format conversion, filtering, or creating encrypted copies.

## Prerequisites

- A source file in JSON or YAML format containing key-value pairs

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
| `destination.type` | Must be `"file"` |
| `destination.path` | Output file path (e.g. `"./output/secrets-copy.json"`) |
| `destination.format` | Output format: `"json"` or `"yaml"` |
| `destination.encrypt` | Set to `true` to encrypt the output file |
| `destination.encrypt_key` | Encryption key (required if `encrypt` is `true`) |
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
| `SECRET_SHIFT_DST_PATH` | `destination.path` |
| `SECRET_SHIFT_DST_FORMAT` | `destination.format` |
| `SECRET_SHIFT_DST_ENCRYPT` | `destination.encrypt` |
| `SECRET_SHIFT_DST_ENCRYPT_KEY` | `destination.encrypt_key` |
| `SECRET_SHIFT_DRY_RUN` | `dry_run` |
