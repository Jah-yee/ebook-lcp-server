#!/usr/bin/env sh
set -eu

KUSTOMIZE="${KUSTOMIZE:-kustomize}"

go test ./...
go vet ./...
go build -buildvcs=false ./...

(
  cd frontend
  if [ ! -d node_modules ]; then
    npm ci
  fi
  npm run build
)

yamllint -c .yamllint.yml \
  deploy/k8s \
  deploy/argocd \
  docker-compose.yml \
  deploy/monitoring/prometheus.yml \
  .github/workflows/go.yml \
  .gitlab-ci.yml

for env in dev staging prod; do
  "${KUSTOMIZE}" build "deploy/overlays/${env}" >/dev/null
done

"${KUSTOMIZE}" build deploy/argocd/apps >/dev/null
