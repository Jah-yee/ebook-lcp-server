# Production Deployment Guide

## Security Hardening

### 1. Transport Security

- Always put Traefik / Caddy / Nginx in front
- Force TLS 1.3
- Use Let's Encrypt

### 2. Secrets Management

- Use Docker secrets or Kubernetes Secrets
- Recommended: External secrets operator + Bitwarden / Vault

### 3. Database & Storage

- Switch from JSON/in-memory to PostgreSQL + S3/MinIO
- Enable connection encryption

### 4. Runtime Security

- Non-root user (already done)
- Read-only filesystem where possible
- Resource limits in Docker/K8s

## Monitoring Stack

- Prometheus + Grafana (provided)
- Loki for logs
- Alertmanager rules (example provided)

## Backup Strategy

...

## High Availability

...

## Recommended Production Architecture Diagram (add one)
