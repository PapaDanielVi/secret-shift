# Contributing to SecretShift

Thank you for your interest in contributing!

## Development Setup

1. Clone the repository
2. Install Go 1.26+
3. Run `make test` to verify everything works
4. Run `make lint` to check code style

## Building

```bash
make build
```

The binary will be at `bin/secret-shift`.

## Testing

```bash
make test          # Run all tests with race detector
make cover         # Generate coverage report
```

## Linting

```bash
make lint
```

Uses `golangci-lint` with the project's `.golangci.yml` configuration.

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Ensure `make test` and `make lint` pass
5. Open a PR with a clear description of the changes

## Code Style

- Follow standard Go conventions
- Keep changes minimal and focused
- Don't refactor unrelated code
- Use `make lint` to catch style issues automatically
