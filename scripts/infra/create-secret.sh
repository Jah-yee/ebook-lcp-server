#!/usr/bin/env sh
set -eu

if [ "$#" -ne 2 ]; then
  echo "usage: scripts/infra/create-secret.sh <namespace> <env-file>" >&2
  exit 1
fi

namespace="$1"
env_file="$2"

kubectl create namespace "${namespace}" --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic lcp-secrets \
  --namespace "${namespace}" \
  --from-literal=jwt-secret="$(grep '^JWT_SECRET=' "${env_file}" | cut -d= -f2-)" \
  --from-literal=admin-2fa-code="$(grep '^ADMIN_2FA_CODE=' "${env_file}" | cut -d= -f2-)" \
  --from-literal=postgres-user="$(grep '^POSTGRES_USER=' "${env_file}" | cut -d= -f2-)" \
  --from-literal=postgres-password="$(grep '^POSTGRES_PASSWORD=' "${env_file}" | cut -d= -f2-)" \
  --from-literal=db-dsn="$(grep '^DB_DSN=' "${env_file}" | cut -d= -f2-)" \
  --from-literal=monitoring-token="$(grep '^MONITORING_TOKEN=' "${env_file}" | cut -d= -f2-)" \
  --dry-run=client -o yaml | kubectl apply -f -
