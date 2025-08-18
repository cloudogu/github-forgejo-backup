#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

cd "$(dirname "$(realpath "$0")")/.."

SPEC_URL='https://forgejo.cloudogu.com/swagger.v1.json'
GEN_IMAGE='openapitools/openapi-generator-cli:v7.14.0'

PKG_NAME='forgejo'
PKG_PATH="internal/${PKG_NAME}"

rm -rf "${PKG_PATH}"
mkdir -p "${PKG_PATH}"

curl -fsSL -o "${PKG_PATH}/swagger.yaml" "${SPEC_URL}"

docker run --rm \
    -u "$(id -u):$(id -g)" \
    -v "${PWD}:/local" \
    "${GEN_IMAGE}" generate \
    -i /local/internal/forgejo/swagger.yaml \
    -g go \
    -o /local/internal/${PKG_NAME} \
    --package-name "${PKG_NAME}" \
    --global-property=apis,models,apiDocs=false,modelDocs=false,apiTests=false,modelTests=false \
    --additional-properties=packageName=${PKG_NAME},isGoSubmodule=true,packageVersion=0.1.0,withGoCodegenComment=true,enumClassPrefix=true,structPrefix=true

git add internal/forgejo/
git commit \
    -m "Updated Forgejo API Client" \
    -m "From Spec: ${SPEC_URL}" \
    -m "With Tool: ${GEN_IMAGE}"
