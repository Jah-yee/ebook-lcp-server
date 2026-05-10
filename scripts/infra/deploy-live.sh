#!/usr/bin/env sh
set -eu

if [ -z "${KUBECONFIG:-}" ]; then
  echo "KUBECONFIG must point to a kubeconfig for the K3s cluster" >&2
  exit 1
fi

KUSTOMIZE="${KUSTOMIZE:-kustomize}"
REGISTRY_NAMESPACE="${REGISTRY_NAMESPACE:-registry}"
ARGOCD_NAMESPACE="${ARGOCD_NAMESPACE:-argocd}"

kubectl create namespace "${REGISTRY_NAMESPACE}" --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace "${ARGOCD_NAMESPACE}" --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace lcp-dev --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace lcp-staging --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace lcp-prod --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -k deploy/registry
kubectl apply -n "${ARGOCD_NAMESPACE}" -f deploy/argocd/root-application.yaml

kubectl apply -k deploy/overlays/dev
kubectl apply -k deploy/overlays/staging
kubectl apply -k deploy/overlays/prod

cat <<'EOF'
Live deployment manifests have been applied.

Next:
1. Install cert-manager and create the letsencrypt-prod ClusterIssuer.
2. Install or confirm Traefik ingress in K3s.
3. Wait for ArgoCD to reconcile the applications.
4. Check the rollout status of lcp-core, lsd-core, lcp-server, and lcp-admin-ui.
EOF
