#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="${1:-round-robin-raw}"
kubectl -n "${NAMESPACE}" set env deployment/round-robin-app FAIL_READY-

