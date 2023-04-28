#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if ! which kubectl >/dev/null; then
  echo "Error: kubectl not found in path"
  exit 1
fi

cd $(dirname "${BASH_SOURCE[0]}")/..

kubectl apply -f ./crds/clustersecrets.yaml

HOST=host.internal CACERT=$(cat ./.local/ssl/ca.pem | openssl base64 -A) envsubst < ./.local/k8s-resources.yaml | kubectl apply -f -
