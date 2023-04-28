#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if ! which go >/dev/null; then
  echo "Error: go not found in path"
  exit 1
fi

export GOBIN=$(realpath $(dirname "${BASH_SOURCE[0]}")/../bin)

go install github.tools.sap/cs-devops/gotpl/cmd/gotpl
go install golang.org/x/tools/cmd/goimports
