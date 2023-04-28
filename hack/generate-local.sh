#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if ! which cfssl >/dev/null; then
  echo "Error: cfssl not found in path"
  exit 1
fi

if ! which cfssljson >/dev/null; then
  echo "Error: cfssljson not found in path"
  exit 1
fi

cd $(dirname "${BASH_SOURCE[0]}")/../.local/ssl

cfssl gencert -initca ca.json | cfssljson -bare ca
cfssl gencert -ca ca.pem -ca-key ca-key.pem webhook.json  | cfssljson -bare webhook
