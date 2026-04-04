#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
ENV_FILE="${ROOT_DIR}/.env"

if [[ -f "${ENV_FILE}" ]]; then
  # shellcheck disable=SC1090
  source "${ENV_FILE}"
fi

CLUSTER_NAME="${KIND_CLUSTER_NAME:-round-robin}"
APP_IMAGE="${APP_IMAGE:-round-robin-app:local}"
WORKER_IMAGE="${WORKER_IMAGE:-round-robin-worker:local}"

kind load docker-image --name "${CLUSTER_NAME}" "${APP_IMAGE}" "${WORKER_IMAGE}"

