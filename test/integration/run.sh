#!/usr/bin/env bash
set -euxo pipefail

KUBECONFIG=/var/snap/microk8s/current/credentials/client.config
REGISTRY="localhost:32000"
IMG="$REGISTRY/capi-bootstrap-charmed-k8s-controller:latest"

cd "$(dirname "$(realpath $0)")/../.."
export KUBECONFIG="$KUBECONFIG"
make docker-build docker-push deploy "IMG=$IMG" "DEPLOY_IMG=$IMG"
kubectl rollout status -n capi-charmed-k8s-bootstrap-system deployment/capi-charmed-k8s-bootstrap-controller-manager --timeout 1m
export USE_EXISTING_CLUSTER=true
make test
