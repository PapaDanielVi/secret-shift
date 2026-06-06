# Context Glossary

## Core Concepts

| Term | Definition |
|------|------------|
| **Secret** | A key/value pair with optional type metadata. Represents an environment variable or secret value to be migrated between providers. |
| **Secret Type** | One of `env`, `secret`, or `file`. Determines how a secret is written to destinations. |
| **Provider** | A source or destination for secrets (GitHub, GitLab, Vault, etcd, Kubernetes, File). |
| **Source** | A Provider that reads secrets during sync. |
| **Destination** | A Provider that writes secrets during sync. |
| **Pipeline** | A sync operation: Source → Process → Destination. |
| **Sync** | One execution of the pipeline. |
| **Conflict Strategy** | One of `replace`, `skip`, or `report` — how to handle existing secrets at destinations. |

## Secret Types

| Type | Behavior |
|------|----------|
| **env** | Environment variable. For GitHub: Actions variables (plain text). For GitLab: CI/CD variables. For Kubernetes: ConfigMaps. |
| **secret** | Encrypted secret. For GitHub: Actions secrets (RSA-OAEP encrypted). For Kubernetes: Secrets. |
| **file** | File-type variable. For GitLab: CI/CD file variables (mounted as files in runners). Currently not used as input but preserved in GitLab-to-GitLab sync. |

## Providers

| Provider | Source Capability | Destination Capability | Notes |
|----------|-------------------|---------------------|-------|
| **GitHub** | Actions secrets + variables | Actions secrets + variables | Secrets require RSA-OAEP encryption. Variables are plain text. |
| **GitLab** | CI/CD variables (env_var + file) | CI/CD variables (env_var + file) | File variables are mounted in CI/CD runners. |
| **Vault** | KV v2 secrets | KV v2 secrets | Single KV entry per secret. |
| **etcd** | Key-value pairs under prefix | Key-value store | One key per secret. |
| **Kubernetes** | Secrets + ConfigMaps | Secrets + ConfigMaps | Routes by secret type. |
| **File** | JSON/YAML key/value | JSON/YAML file | AES-256-GCM encryption optional. All values typed as `env`. |

## Processing

| Term | Definition |
|------|------------|
| **Add Prefix** | Prepends a string to all secret names during processing. |
| **Add Suffix** | Appends a string to all secret names during processing. |
| **Include Regex** | Filters secrets by name pattern (include only matches). |
| **Exclude Regex** | Filters secrets by name pattern (exclude matches). |
| **Include Types** | Filters by secret type (env, secret, file). |
| **Exclude Types** | Excludes secrets by type. |

## Run Modes

| Mode | Description |
|------|-------------|
| **One-shot** | Single sync execution. |
| **Periodic** | Continuous sync at a fixed interval. |
| **Cron** | Scheduled sync via cron expression. |
| **Server** | HTTP server with health endpoints + periodic sync. |

## Environment Variable Conventions

| Variable | Purpose |
|----------|---------|
| `SECRET_SHIFT_SRC_TYPE` | Source provider type. |
| `SECRET_SHIFT_DST_TYPE` | Destination provider type. |
| `SECRET_SHIFT_SRC_<PROVIDER>_TOKEN` | Token for source provider (e.g., `SECRET_SHIFT_SRC_GITHUB_TOKEN`). |
| `SECRET_SHIFT_DST_<PROVIDER>_TOKEN` | Token for destination provider. |
| `SECRET_SHIFT_DRY_RUN` | Enable dry-run mode. |

Token resolution order: direct config value → `token_env` field → `SECRET_SHIFT_SRC/DST_<PROVIDER>_TOKEN`.