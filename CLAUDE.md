# Secret-Shift Project Glossary & Conventions

See [CONTEXT.md](./CONTEXT.md) for domain terminology.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

## 5. Provider Patterns

**Context threading convention:**

All provider constructors that perform network initialization accept `context.Context` as their first parameter:

```go
func New(ctx context.Context, ...) (*Provider, error)
```

This applies to:
- `provider/github.New()` — passes context to oauth2 HTTP client
- `provider/gitlab.New()` — retained for interface consistency

Providers that don't need network initialization during construction (vault, etcd, kubernetes, file) do not take a context in `New()`. Context is passed at the method level (`Read(ctx)`, `Write(ctx)`) for all providers.

**Registry pattern:**

All providers register themselves via `init()` in `internal/provider/registry.go`. Each provider provides a `SourceFactory` and `DestFactory` function. The pipeline uses `provider.CreateSource()` and `provider.CreateDestination()` to construct providers from config, passing options as `map[string]any`.

When adding a new provider:
1. Implement `Source` and/or `Destination` interfaces
2. Add `init()` that calls `provider.Register()`
3. Add a `getString()` helper in the provider package for extracting string options
4. Update `pipeline.go` `buildSourceOpts()` and `buildDestOpts()` to map config fields

**Per-step environment variables:**

Source and destination each resolve tokens independently. Resolution order:
1. Direct config value (`token` field)
2. `token_env` field (custom env var name)
3. `SECRET_SHIFT_SRC_<PROVIDER>_TOKEN` or `SECRET_SHIFT_DST_<PROVIDER>_TOKEN`

The suffix is uppercase provider name: `GITHUB`, `GITLAB`, `VAULT`, `ETCD`, `KUBERNETES`, `FILE`.

## 6. File Provider

The file provider (`internal/provider/file/`) implements both `Source` and `Destination`. It supports:
- JSON and YAML formats
- AES-256-GCM encryption (same encrypt/decrypt for both read and write)
- Default format is JSON if not specified

## 7. Server Mode

The health server (`internal/server/health.go`) uses only the standard library (`net/http`). It tracks sync success/error with mutex-protected timestamps and atomic counters. Endpoints:
- `/healthz` — always 200 (liveness)
- `/readyz` — 200 after first successful sync, 503 otherwise (readiness)
- `/status` — JSON with sync count, error count, timestamps

## 8. Testing Conventions

- Registry tests live in `internal/provider/registry_test.go` as `package provider_test` with blank imports (`_ "..."`) to trigger `init()` functions
- Health server tests use `httptest.NewRecorder` for handler tests and an actual HTTP server on port 18080 for integration tests
- Provider tests use compile-time interface checks: `var _ provider.Source = (*Provider)(nil)`
- `t.TempDir()` for file-based tests
