#!/usr/bin/env sh
set -eu

if [ -z "${KUBECONFIG:-}" ]; then
  echo "KUBECONFIG must point to a kubeconfig for the K3s cluster" >&2
  exit 1
fi

NAMESPACE_ARGOCD="${NAMESPACE_ARGOCD:-argocd}"
NAMESPACE_CERT_MANAGER="${NAMESPACE_CERT_MANAGER:-cert-manager}"
NAMESPACE_MONITORING="${NAMESPACE_MONITORING:-monitoring}"

kubectl create namespace "${NAMESPACE_ARGOCD}" --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace "${NAMESPACE_CERT_MANAGER}" --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace "${NAMESPACE_MONITORING}" --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace registry --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace lcp-dev --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace lcp-staging --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace lcp-prod --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -k deploy/registry
kubectl apply -n "${NAMESPACE_ARGOCD}" -f deploy/argocd/root-application.yaml

cat <<'EOF'
Cluster bootstrap complete at the namespace level.

Next install steps on the K3s cluster:
1. ArgoCD
2. cert-manager
3. Prometheus Operator stack
4. optional external secrets or Vault

Then sync the ArgoCD root application.
EOF
