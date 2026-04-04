#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="${1:-round-robin}"
DEPLOYMENT="${2:-round-robin-app}"
kubectl -n "${NAMESPACE}" rollout undo deployment/"${DEPLOYMENT}"

