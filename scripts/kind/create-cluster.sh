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
REGISTRY_PORT="${KIND_REGISTRY_PORT:-5001}"
NODE_IMAGE="${KIND_NODE_IMAGE:-kindest/node:v1.31.0@sha256:53df588e04085fd41ae12de0c3fe4c72f7013bba32a20e7325357a1ac94ba865}"

if ! docker inspect "${REGISTRY_NAME}" >/dev/null 2>&1; then
  docker run -d --restart=always -p "127.0.0.1:${REGISTRY_PORT}:5000" --name "${REGISTRY_NAME}" registry:2
fi

if kind get clusters | grep -qx "${CLUSTER_NAME}"; then
  echo "kind cluster ${CLUSTER_NAME} already exists"
else
  CONFIG_FILE="$(mktemp)"
  trap 'rm -f "${CONFIG_FILE}"' EXIT

  cat >"${CONFIG_FILE}" <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${REGISTRY_PORT}"]
    endpoint = ["http://${REGISTRY_NAME}:5000"]
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
EOF

  kind create cluster --name "${CLUSTER_NAME}" --image "${NODE_IMAGE}" --config "${CONFIG_FILE}"
fi

kind export kubeconfig --name "${CLUSTER_NAME}"
kubectl wait --for=condition=Ready node --all --timeout=180s

docker network connect kind "${REGISTRY_NAME}" >/dev/null 2>&1 || true

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${REGISTRY_PORT}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF

kubectl cluster-info --context "kind-${CLUSTER_NAME}"
