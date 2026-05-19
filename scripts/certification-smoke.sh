#!/bin/sh
set -eu

base_url="${BASE_URL:-http://localhost:8080}"
out="${1:-certification-report.json}"

token="$(
  curl -fsS "$base_url/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin","twoFactor":"123456"}' |
    jq -r .token
)"

health="$(curl -fsS "$base_url/healthz")"
ready="$(curl -fsS "$base_url/readyz")"
status="$(curl -fsS "$base_url/api/v1/lcp/status" -H "Authorization: Bearer $token")"

jq -n \
  --arg generatedAt "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  --argjson health "$health" \
  --argjson ready "$ready" \
  --argjson status "$status" \
  '{generatedAt: $generatedAt, checks: {healthz: $health, readyz: $ready, lcpStatus: $status}}' > "$out"

echo "$out"
