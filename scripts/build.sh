#!/usr/bin/env bash

set -euxo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(cd "${script_dir}/.." && pwd)"
pushd "${root_dir}" >/dev/null || exit 1

subdivision=${SUBDIVISION:-zan}

# support other container tools, e.g. podman
CONTAINER_CLI=${CONTAINER_CLI:-docker}
# Use buildx for CI by default, allow overriding for old clients or other tools like podman
CONTAINER_BUILDER=${CONTAINER_BUILDER:-"buildx build"}

PUSH=""
if [ -z "${DRY_RUN:-}" ]; then
  PUSH="--push"
fi

HUB=${HUB:-"docker.io/adhp"}
TAG=${TAG:-"dirty"}
TAG="${subdivision}"-"${TAG}"

if [[ -f "frontend/config.json" ]]; then
  mv frontend/config.json frontend/config.json.bak
fi

cp "frontend/configs/${subdivision}.json" frontend/config.json

${CONTAINER_CLI} ${CONTAINER_BUILDER} \
  --target go_final \
  ${PUSH} \
  --tag "${HUB}/ids-backend:${TAG}" \
  .

${CONTAINER_CLI} ${CONTAINER_BUILDER} \
  --target node_final \
  ${PUSH} \
  --tag "${HUB}/ids-frontend:${TAG}" \
  .

if [[ -f "frontend/config.json.bak" ]]; then
  mv frontend/config.json.bak frontend/config.json
fi
