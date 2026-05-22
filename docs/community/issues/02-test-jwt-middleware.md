# Add focused tests for JWT middleware behavior

Labels: `good first issue`, `testing`, `help wanted`

## Summary

Add unit tests around the JWT middleware in `internal/adapter/jwt` so auth failures and role checks are covered more explicitly.

## Why this matters

Authentication behavior is high-impact, and the middleware package currently has no direct test coverage.

## Suggested files

- `internal/adapter/jwt/middleware.go`
- `internal/adapter/jwt/middleware_test.go`

## Acceptance criteria

- Tests cover a valid token request.
- Tests cover missing token and malformed token cases.
- Tests cover a role mismatch returning an authorization failure.
- `go test ./...` passes.

## Nice to have

- Use table-driven tests to keep future cases easy to add.
