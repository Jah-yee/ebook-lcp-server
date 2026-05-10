# Support and Knowledge Transfer

## Warranty Support

The signed contract includes 3 months of basic support after final delivery for technical bugs and defects.

## Critical Incidents

Critical issues require emergency response within 2 hours. A critical issue is any defect that makes the main LCP service unavailable or prevents the core processing API from working.

Recommended production alert routes:

- API 5xx rate above agreed threshold.
- Kubernetes Deployment unavailable.
- PostgreSQL unavailable.
- Persistent volume nearing capacity.
- TLS certificate close to expiry.

## Knowledge Transfer

The contract requires at least 8 hours of training for the employer's internal team. Suggested agenda:

1. Architecture and data flow.
2. Local development and Docker Compose.
3. Authentication, roles, and admin 2FA.
4. REST and GraphQL API usage.
5. PostgreSQL migrations and backup recovery.
6. Kubernetes deployment, HPA, Ingress, Secrets, and ConfigMaps.
7. Monitoring, logs, alerts, and incident response.
8. CI/CD operation and release process.

## Handover Checklist

- Repository access granted.
- Production Secrets rotated by employer.
- DNS and TLS configured for the final domain.
- `kubectl apply -k deploy/k8s` succeeds in the target namespace.
- `go test ./...`, `go vet ./...`, and frontend build pass in CI.
- Load test report attached for acceptance testing.
- Trivy scan shows no critical vulnerabilities.
