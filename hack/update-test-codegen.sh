#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if ! which go >/dev/null; then
  echo "Error: go not found in path"
  exit 1
fi

cd $(dirname "${BASH_SOURCE[0]}")/..

CODEGEN_PKG=$(go list -m -f '{{.Dir}}' github.tools.sap/cs-devops/kubernetes-testing)
OUTPUT_BASE=$(mktemp -d)
trap 'rm -rf "${OUTPUT_BASE}"' EXIT

# echo "PWD: ${PWD}"
# echo "CODEGEN_PKG: ${CODEGEN_PKG}"
# echo "OUTPUT_BASE: ${OUTPUT_BASE}"

bin/gotpl -f ./test/generate.json "${CODEGEN_PKG}"/templates/environment.tpl > "${OUTPUT_BASE}"/zz_environment.go
bin/goimports -w "${OUTPUT_BASE}"/zz_environment.go
cp -f "${OUTPUT_BASE}"/zz_environment.go ./test/zz_environment.go
