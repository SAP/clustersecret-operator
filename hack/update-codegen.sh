#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if ! which go >/dev/null; then
  echo "Error: go not found in path"
  exit 1
fi

cd "$(dirname "${BASH_SOURCE[0]}")"/..

go mod download k8s.io/code-generator
CODEGEN_PKG=$(go list -m -f '{{.Dir}}' k8s.io/code-generator)
GEN_PKG_PATH=$(go list -m)/pkg
OUTPUT_PATH=$(mktemp -d)
trap 'rm -rf "${OUTPUT_PATH}"' EXIT

# echo "PWD: ${PWD}"
# echo "CODEGEN_PKG: ${CODEGEN_PKG}"
# echo "GEN_PKG_PATH: ${GEN_PKG_PATH}"
# echo "OUTPUT_PATH: ${OUTPUT_PATH}"

source "${CODEGEN_PKG}"/kube_codegen.sh

kube::codegen::gen_helpers \
  --boilerplate ./hack/boilerplate.go.txt \
  ./pkg/apis

kube::codegen::gen_client \
  --with-watch \
  --with-applyconfig \
  --output-dir "${OUTPUT_PATH}"/client \
  --output-pkg "${GEN_PKG_PATH}"/client \
  --boilerplate ./hack/boilerplate.go.txt \
  ./pkg/apis

rm -rf "./pkg/client" && cp -Rf "${OUTPUT_PATH}"/client ./pkg
