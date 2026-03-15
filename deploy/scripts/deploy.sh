#!/usr/bin/env bash
set -euo pipefail

# Deploy script for todo-app on zima-lab cluster.
# Usage: ./deploy/scripts/deploy.sh [--image IMAGE] [--namespace NAMESPACE]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
DEPLOY_DIR="$PROJECT_ROOT/deploy"

IMAGE="${IMAGE:-forgia-example-todo:latest}"
NAMESPACE="${NAMESPACE:-todo-app}"

# --- Helpers ---

deploy_log() {
  echo "[deploy] $*"
}

deploy_err() {
  echo "[deploy] ERROR: $*" >&2
}

# --- Prerequisites ---

deploy_check_prerequisites() {
  local missing=0

  for tool in kubectl docker; do
    if ! command -v "$tool" >/dev/null 2>&1; then
      deploy_err "'$tool' is required but not found in PATH"
      missing=1
    fi
  done

  if [[ "$missing" -ne 0 ]]; then
    exit 1
  fi

  # Verify cluster connectivity
  if ! kubectl cluster-info >/dev/null 2>&1; then
    deploy_err "cannot connect to Kubernetes cluster — check your kubeconfig"
    exit 1
  fi

  deploy_log "prerequisites OK"
}

# --- Parse args ---

deploy_parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --image)
        IMAGE="${2:?--image requires a value}"
        shift 2
        ;;
      --namespace)
        NAMESPACE="${2:?--namespace requires a value}"
        shift 2
        ;;
      --help)
        echo "Usage: $0 [--image IMAGE] [--namespace NAMESPACE]"
        echo ""
        echo "Options:"
        echo "  --image IMAGE         Docker image for the operator (default: forgia-example-todo:latest)"
        echo "  --namespace NAMESPACE Kubernetes namespace (default: todo-app)"
        exit 0
        ;;
      *)
        deploy_err "unknown argument: $1"
        exit 64
        ;;
    esac
  done
}

# --- Build ---

deploy_build_image() {
  deploy_log "building operator image: $IMAGE"
  docker build -t "$IMAGE" -f "$PROJECT_ROOT/cmd/operator/Dockerfile" "$PROJECT_ROOT"
}

# --- Deploy ---

deploy_apply_manifests() {
  deploy_log "applying CRD"
  kubectl apply -f "$DEPLOY_DIR/crd/todo-crd.yaml"

  deploy_log "applying operator manifests to namespace=$NAMESPACE"
  kubectl apply -f "$DEPLOY_DIR/operator/namespace.yaml"
  kubectl apply -f "$DEPLOY_DIR/operator/serviceaccount.yaml"
  kubectl apply -f "$DEPLOY_DIR/operator/rbac.yaml"
  kubectl apply -f "$DEPLOY_DIR/operator/resource-quota.yaml"

  # Patch the image in the deployment if a custom image is provided
  kubectl apply -f "$DEPLOY_DIR/operator/deployment.yaml"
  if [[ "$IMAGE" != "forgia-example-todo:latest" ]]; then
    kubectl set image deployment/todo-operator \
      operator="$IMAGE" \
      -n "$NAMESPACE"
  fi

  kubectl apply -f "$DEPLOY_DIR/operator/service.yaml"

  deploy_log "applying plugin provisioning"
  kubectl apply -f "$DEPLOY_DIR/plugin/grafana-plugin-provisioning.yaml" 2>/dev/null || \
    deploy_log "plugin provisioning is a local config file — copy it to Grafana provisioning directory manually"
}

# --- Verify ---

deploy_verify() {
  deploy_log "waiting for operator pod to be ready (timeout 120s)"
  if kubectl wait --for=condition=available deployment/todo-operator \
    -n "$NAMESPACE" --timeout=120s; then
    deploy_log "operator deployment is available"
  else
    deploy_err "operator deployment did not become available within 120s"
    kubectl get pods -n "$NAMESPACE" -l app.kubernetes.io/name=todo-operator
    exit 1
  fi
}

# --- Main ---

main() {
  deploy_parse_args "$@"
  deploy_check_prerequisites
  deploy_build_image
  deploy_apply_manifests
  deploy_verify
  deploy_log "deployment complete"
}

main "$@"
