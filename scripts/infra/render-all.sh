#!/usr/bin/env sh
set -eu

KUSTOMIZE="${KUSTOMIZE:-kustomize}"

for env in dev staging prod; do
  echo "Rendering ${env}"
  "${KUSTOMIZE}" build "deploy/overlays/${env}" > "/tmp/lcp-${env}.yaml"
  echo "Wrote /tmp/lcp-${env}.yaml"
done

echo "Rendering ArgoCD app-of-apps"
"${KUSTOMIZE}" build deploy/argocd/apps > /tmp/lcp-argocd-apps.yaml
echo "Wrote /tmp/lcp-argocd-apps.yaml"
