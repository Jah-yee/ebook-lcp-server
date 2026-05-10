# Architecture

## Components

- HTTP server: standard-library `net/http` router in `cmd/server`.
- Auth middleware: HS256 JWT validation and RBAC in `internal/auth`.
- REST adapter: contract endpoints in `internal/adapter/rest`.
- GraphQL adapter: legacy publication and license operations in `internal/adapter/graphql`.
- Use cases: publication upload/encryption and license issuing in `internal/usecase/lcp`.
- Repositories: in-memory by default, JSON-backed metadata when `DATA_DIR` is configured.
- Storage: filesystem-backed encrypted content under `LCP_STORAGE_FS_DIR`.

## Data Flow

```text
Client
  -> JWT/RBAC middleware
  -> REST or GraphQL handler
  -> LCP use case
  -> encrypter/license service
  -> metadata repository
  -> filesystem persistent volume
```

## Production Notes

The local JSON repository is suitable for a single-writer deployment and acceptance testing. For full PostgreSQL replication/sharding requirements, implement the existing repository interfaces with PostgreSQL adapters and run the included migrations in `migrations/`.
