#!/usr/bin/env bash
set -euo pipefail

# Validate Kubernetes YAML manifests for the todo-app deployment.
# Checks syntax, required fields, and security constraints.
# Usage: ./tests/deployment/validate.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
DEPLOY_DIR="$PROJECT_ROOT/deploy"

PASS=0
FAIL=0
ERRORS=""

# --- Helpers ---

validate_log() {
  echo "[validate] $*"
}

validate_err() {
  echo "[validate] FAIL: $*" >&2
  FAIL=$((FAIL + 1))
  ERRORS="${ERRORS}\n  - $*"
}

validate_pass() {
  echo "[validate] PASS: $*"
  PASS=$((PASS + 1))
}

# --- YAML syntax check ---

validate_yaml_syntax() {
  local file="$1"
  local basename
  basename="$(basename "$file")"

  # Try validators in order: python3+pyyaml, kubectl dry-run, basic checks
  if python3 -c "import yaml" 2>/dev/null; then
    if python3 -c "
import yaml, sys
with open(sys.argv[1]) as f:
    list(yaml.safe_load_all(f))
" "$file" 2>/dev/null; then
      validate_pass "YAML syntax OK: $basename"
    else
      validate_err "YAML syntax invalid: $basename"
    fi
  elif command -v kubectl >/dev/null 2>&1 && kubectl cluster-info >/dev/null 2>&1; then
    if kubectl apply --dry-run=client -f "$file" >/dev/null 2>&1; then
      validate_pass "YAML syntax OK (kubectl): $basename"
    else
      validate_err "YAML syntax invalid: $basename"
    fi
  else
    # Basic structural check: verify file is not empty and has valid YAML-like structure
    if [[ -s "$file" ]] && grep -qE "^[a-zA-Z]" "$file"; then
      validate_pass "YAML structure OK (basic): $basename"
    else
      validate_err "YAML appears empty or malformed: $basename"
    fi
  fi
}

# --- Kubeconform validation ---

validate_with_kubeconform() {
  if ! command -v kubeconform >/dev/null 2>&1; then
    validate_log "kubeconform not found — skipping schema validation (install: brew install kubeconform)"
    return
  fi

  local file="$1"
  local basename
  basename="$(basename "$file")"

  if kubeconform -strict -summary "$file" 2>/dev/null; then
    validate_pass "kubeconform schema OK: $basename"
  else
    validate_err "kubeconform schema invalid: $basename"
  fi
}

# --- Security checks ---

validate_no_hardcoded_secrets() {
  local file="$1"
  local basename
  basename="$(basename "$file")"

  # Check for common secret patterns (passwords, tokens, keys with actual values)
  if grep -qEi '(password|token|api_key|secret_key)\s*[:=]\s*["\x27][^<]' "$file" 2>/dev/null; then
    validate_err "potential hardcoded secret in: $basename"
  else
    validate_pass "no hardcoded secrets: $basename"
  fi
}

validate_deployment_security() {
  local file="$DEPLOY_DIR/operator/deployment.yaml"

  # Check runAsNonRoot
  if grep -q "runAsNonRoot: true" "$file"; then
    validate_pass "deployment: runAsNonRoot is set"
  else
    validate_err "deployment: runAsNonRoot not set"
  fi

  # Check allowPrivilegeEscalation
  if grep -q "allowPrivilegeEscalation: false" "$file"; then
    validate_pass "deployment: allowPrivilegeEscalation is false"
  else
    validate_err "deployment: allowPrivilegeEscalation not disabled"
  fi

  # Check readOnlyRootFilesystem
  if grep -q "readOnlyRootFilesystem: true" "$file"; then
    validate_pass "deployment: readOnlyRootFilesystem is set"
  else
    validate_err "deployment: readOnlyRootFilesystem not set"
  fi

  # Check resource limits
  if grep -q "limits:" "$file"; then
    validate_pass "deployment: resource limits defined"
  else
    validate_err "deployment: no resource limits"
  fi
}

# --- RBAC checks ---

validate_rbac_least_privilege() {
  local file="$DEPLOY_DIR/operator/rbac.yaml"

  # Ensure no wildcard permissions
  if grep -q 'verbs:.*"\*"' "$file" 2>/dev/null || grep -q "verbs:.*'\\*'" "$file" 2>/dev/null; then
    validate_err "RBAC: wildcard verbs found (violates least privilege)"
  else
    validate_pass "RBAC: no wildcard verbs"
  fi

  # Ensure only todo.grafana.app apiGroup
  if grep -q "apiGroups:" "$file"; then
    local non_todo_groups
    non_todo_groups=$(grep -A1 "apiGroups:" "$file" | grep -v "apiGroups:" | grep -v "todo.grafana.app" | grep -v "^--$" | grep -v "^\s*$" || true)
    if [[ -z "$non_todo_groups" ]]; then
      validate_pass "RBAC: scoped to todo.grafana.app only"
    else
      validate_err "RBAC: permissions granted outside todo.grafana.app"
    fi
  fi
}

# --- Dockerfile checks ---

validate_dockerfile() {
  local file="$PROJECT_ROOT/cmd/operator/Dockerfile"

  if [[ ! -f "$file" ]]; then
    validate_err "Dockerfile not found at cmd/operator/Dockerfile"
    return
  fi

  # Multi-stage build
  local stage_count
  stage_count=$(grep -ci "^FROM " "$file")
  if [[ "$stage_count" -ge 2 ]]; then
    validate_pass "Dockerfile: multi-stage build ($stage_count stages)"
  else
    validate_err "Dockerfile: not a multi-stage build"
  fi

  # Non-root user
  if grep -q "USER" "$file"; then
    validate_pass "Dockerfile: non-root USER set"
  else
    validate_err "Dockerfile: no USER directive (runs as root)"
  fi
}

# --- Deploy script checks ---

validate_deploy_script() {
  local file="$PROJECT_ROOT/deploy/scripts/deploy.sh"

  if [[ ! -f "$file" ]]; then
    validate_err "deploy script not found at deploy/scripts/deploy.sh"
    return
  fi

  # set -euo pipefail
  if grep -q "set -euo pipefail" "$file"; then
    validate_pass "deploy script: set -euo pipefail present"
  else
    validate_err "deploy script: missing set -euo pipefail"
  fi

  # command -v checks
  if grep -q "command -v" "$file"; then
    validate_pass "deploy script: prerequisite checks with command -v"
  else
    validate_err "deploy script: no prerequisite checks"
  fi

  # shellcheck (if available)
  if command -v shellcheck >/dev/null 2>&1; then
    if shellcheck "$file" 2>/dev/null; then
      validate_pass "deploy script: shellcheck passed"
    else
      validate_err "deploy script: shellcheck found issues"
    fi
  else
    validate_log "shellcheck not found — skipping lint"
  fi
}

# --- Main ---

main() {
  validate_log "starting manifest validation"
  echo ""

  # Validate all YAML files
  for yaml_file in "$DEPLOY_DIR"/crd/*.yaml "$DEPLOY_DIR"/operator/*.yaml "$DEPLOY_DIR"/plugin/*.yaml; do
    if [[ -f "$yaml_file" ]]; then
      validate_yaml_syntax "$yaml_file"
      validate_with_kubeconform "$yaml_file"
      validate_no_hardcoded_secrets "$yaml_file"
    fi
  done

  echo ""

  # Security and RBAC
  validate_deployment_security
  validate_rbac_least_privilege

  echo ""

  # Dockerfile
  validate_dockerfile

  echo ""

  # Deploy script
  validate_deploy_script

  echo ""

  # Summary
  echo "=============================="
  echo "  Results: $PASS passed, $FAIL failed"
  echo "=============================="

  if [[ "$FAIL" -gt 0 ]]; then
    echo -e "\nFailures:$ERRORS"
    exit 1
  fi

  validate_log "all checks passed"
}

main "$@"
