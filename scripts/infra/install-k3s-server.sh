#!/usr/bin/env sh
set -eu

if [ -z "${K3S_TOKEN:-}" ]; then
  echo "K3S_TOKEN must be set" >&2
  exit 1
fi

INSTALL_K3S_EXEC="${INSTALL_K3S_EXEC:-server --cluster-init}"

curl -sfL https://get.k3s.io | \
  INSTALL_K3S_EXEC="${INSTALL_K3S_EXEC} --tls-san ${TLS_SAN:-lcp.example.com}" \
  K3S_TOKEN="${K3S_TOKEN}" \
  sh -
