#!/usr/bin/env sh
set -eu

REGISTRY="${REGISTRY:-registry.testmedical.ir:5000}"
BACKEND_IMAGE="${REGISTRY}/lcp-server"
FRONTEND_IMAGE="${REGISTRY}/lcp-admin-ui"

export BUILDAH_ISOLATION="${BUILDAH_ISOLATION:-chroot}"

buildah bud -t "${BACKEND_IMAGE}:latest" .
buildah bud -t "${FRONTEND_IMAGE}:latest" frontend
buildah push --tls-verify=false "${BACKEND_IMAGE}:latest" "docker://${BACKEND_IMAGE}:latest"
buildah push --tls-verify=false "${FRONTEND_IMAGE}:latest" "docker://${FRONTEND_IMAGE}:latest"
