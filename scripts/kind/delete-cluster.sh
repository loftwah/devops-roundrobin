#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
ENV_FILE="${ROOT_DIR}/.env"

if [[ -f "${ENV_FILE}" ]]; then
  # shellcheck disable=SC1090
  source "${ENV_FILE}"
fi

CLUSTER_NAME="${KIND_CLUSTER_NAME:-round-robin}"
REGISTRY_NAME="${KIND_REGISTRY_NAME:-kind-registry}"
KIND_CONTROL_PLANE_CONTAINER="${CLUSTER_NAME}-control-plane"
kind_missing_with_live_cluster=0

if command -v kind >/dev/null 2>&1; then
  if kind get clusters | grep -qx "${CLUSTER_NAME}"; then
    kind delete cluster --name "${CLUSTER_NAME}"
  else
    echo "kind cluster ${CLUSTER_NAME} does not exist"
  fi
elif docker inspect "${KIND_CONTROL_PLANE_CONTAINER}" >/dev/null 2>&1; then
  echo "kind command not found; cannot delete cluster ${CLUSTER_NAME}" >&2
  kind_missing_with_live_cluster=1
else
  echo "kind command not found; cluster ${CLUSTER_NAME} is already absent"
fi

if docker inspect "${REGISTRY_NAME}" >/dev/null 2>&1; then
  docker rm -f "${REGISTRY_NAME}" >/dev/null
  echo "removed local registry ${REGISTRY_NAME}"
else
  echo "local registry ${REGISTRY_NAME} does not exist"
fi

if [[ "${kind_missing_with_live_cluster}" -eq 1 ]]; then
  exit 1
fi
