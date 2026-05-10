#!/usr/bin/env sh
set -eu

if [ -z "${K3S_TOKEN:-}" ]; then
  echo "K3S_TOKEN must be set" >&2
  exit 1
fi

if [ -z "${K3S_URL:-}" ]; then
  echo "K3S_URL must be set, for example https://server-ip:6443" >&2
  exit 1
fi

curl -sfL https://get.k3s.io | \
  K3S_URL="${K3S_URL}" \
  K3S_TOKEN="${K3S_TOKEN}" \
  sh -
