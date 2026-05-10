# K3s Bootstrap Guide

This project is designed for a self-hosted K3s cluster.

## One-time Host Setup

1. Install K3s server on the first node:

```bash
sudo scripts/infra/install-k3s-config.sh
sudo scripts/infra/install-k3s-registries.sh
sudo K3S_TOKEN=replace-with-token scripts/infra/install-k3s-server.sh
```

2. Join additional workers:

```bash
sudo K3S_URL=https://server-ip:6443 K3S_TOKEN=replace-with-token scripts/infra/install-k3s-agent.sh
```

3. Confirm the cluster:

```bash
kubectl get nodes
```

## Bootstrap Cluster Services

The repo assumes these services are installed in the cluster:

- ArgoCD
- cert-manager
- Prometheus Operator stack
- optional Vault or External Secrets
- self-hosted registry from `deploy/registry`

You can then run:

```bash
KUBECONFIG=~/.kube/config scripts/infra/bootstrap-k3s-stack.sh
```

To publish images into the registry:

```bash
scripts/infra/build-and-push-images.sh
```

## GitOps Flow

1. Apply the ArgoCD root application.
2. ArgoCD creates the `dev`, `staging`, and `prod` Applications.
3. Each Application syncs its own overlay.

The self-hosted K3s cluster uses the built-in Traefik ingress controller and local-path storage by default.
