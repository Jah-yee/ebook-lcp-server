#!/usr/bin/env sh
set -eu

REGISTRY="${REGISTRY:-registry.testmedical.ir:5000}"
TAG="${TAG:-latest}"
BACKEND_IMAGE="${REGISTRY}/lcp-server"
FRONTEND_IMAGE="${REGISTRY}/lcp-admin-ui"

export BUILDAH_ISOLATION="${BUILDAH_ISOLATION:-chroot}"

buildah bud -t "${BACKEND_IMAGE}:${TAG}" .
buildah bud -t "${FRONTEND_IMAGE}:${TAG}" frontend
buildah push --tls-verify=false "${BACKEND_IMAGE}:${TAG}" "docker://${BACKEND_IMAGE}:${TAG}"
buildah push --tls-verify=false "${FRONTEND_IMAGE}:${TAG}" "docker://${FRONTEND_IMAGE}:${TAG}"
