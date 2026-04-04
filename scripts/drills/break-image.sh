#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="${1:-round-robin}"
DEPLOYMENT="${2:-round-robin-app}"
CONTAINER="${3:-app}"
kubectl -n "${NAMESPACE}" set image deployment/"${DEPLOYMENT}" "${CONTAINER}"=ghcr.io/example/does-not-exist:broken

