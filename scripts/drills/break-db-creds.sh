#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="${1:-round-robin}"
DEPLOYMENT="${2:-round-robin-app}"
kubectl -n "${NAMESPACE}" set env deployment/"${DEPLOYMENT}" POSTGRES_PASSWORD=wrong-password

