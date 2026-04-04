#!/usr/bin/env bash
set -euo pipefail

required_commands=(
  docker
  kind
  kubectl
  helm
  terraform
  jq
  yq
  trivy
)

missing=()
for command in "${required_commands[@]}"; do
  if ! command -v "${command}" >/dev/null 2>&1; then
    missing+=("${command}")
  fi
done

if (( ${#missing[@]} > 0 )); then
  printf 'Missing commands: %s\n' "${missing[*]}"
  printf 'Run `nix develop` to enter the pinned toolchain shell.\n'
  exit 1
fi

printf 'Local prerequisites are available.\n'

