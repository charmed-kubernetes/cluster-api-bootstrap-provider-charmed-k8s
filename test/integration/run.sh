#!/usr/bin/env bash
set -euxo pipefail

REGISTRY="localhost:32000"
IMG="$REGISTRY/capi-bootstrap-charmed-k8s-controller:latest"

cd "$(dirname "$(realpath $0)")/../.."
make docker-build docker-push deploy "IMG=$IMG" "DEPLOY_IMG=$IMG"
kubectl rollout status -n capi-charmed-k8s-bootstrap-system deployment/capi-charmed-k8s-bootstrap-controller-manager --timeout 1m
export USE_EXISTING_CLUSTER=true
make test
