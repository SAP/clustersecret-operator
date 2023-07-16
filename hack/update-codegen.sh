#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if ! which go >/dev/null; then
  echo "Error: go not found in path"
  exit 1
fi

cd $(dirname "${BASH_SOURCE[0]}")/..

go mod download k8s.io/code-generator
CODEGEN_PKG=$(go list -m -f '{{.Dir}}' k8s.io/code-generator)
GEN_PKG_PATH=$(go list -m)/pkg
OUTPUT_BASE=$(mktemp -d)

trap 'rm -rf "${OUTPUT_BASE}"' EXIT

# echo "PWD: ${PWD}"
# echo "CODEGEN_PKG: ${CODEGEN_PKG}"
# echo "GEN_PKG_PATH: ${GEN_PKG_PATH}"
# echo "OUTPUT_BASE: ${OUTPUT_BASE}"

/bin/bash "${CODEGEN_PKG}"/generate-groups.sh all \
  "${GEN_PKG_PATH}"/client "${GEN_PKG_PATH}"/apis \
  core.cs.sap.com:v1alpha1 \
  --output-base "${OUTPUT_BASE}"/ \
  --go-header-file ./hack/boilerplate.go.txt

rm -rf "./pkg/client" && cp -Rf "${OUTPUT_BASE}"/"${GEN_PKG_PATH}" .
