#!/usr/bin/env sh
set -eu

SOURCE="${1:-deploy/k3s/registries.yaml.example}"
TARGET="${2:-/etc/rancher/k3s/registries.yaml}"

if [ ! -f "${SOURCE}" ]; then
  echo "registry config source not found: ${SOURCE}" >&2
  exit 1
fi

if [ "$(id -u)" -eq 0 ]; then
  INSTALL_PREFIX=""
else
  INSTALL_PREFIX="sudo"
fi

${INSTALL_PREFIX} install -d -m 0755 "$(dirname "${TARGET}")"
${INSTALL_PREFIX} install -m 0644 "${SOURCE}" "${TARGET}"

echo "installed ${SOURCE} to ${TARGET}"
