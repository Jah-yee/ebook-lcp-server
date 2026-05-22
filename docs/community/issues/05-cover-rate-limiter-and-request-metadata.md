# Add tests for rate limiting and request metadata helpers

Labels: `good first issue`, `testing`, `help wanted`

## Summary

Add tests for the small utility packages behind request rate limiting and request metadata propagation.

## Why this matters

These packages are small, isolated, and a good place for a first contribution that still improves reliability.

## Suggested files

- `internal/ratelimit/limiter.go`
- `internal/requestmeta/requestmeta.go`
- `internal/ratelimit/limiter_test.go`
- `internal/requestmeta/requestmeta_test.go`

## Acceptance criteria

- Tests cover the basic allow/deny flow for the limiter.
- Tests cover storing and retrieving request metadata from context or middleware helpers.
- The new tests are deterministic and do not rely on long sleeps.
- `go test ./...` passes.

## Nice to have

- Use clear test names so the packages stay approachable for future contributors.
