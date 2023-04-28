#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

cd "$(dirname "$0")"

function download() {
  local name=$1
  local tag=$2
  local sources=$3

  rm -rf $name.tmp
  git clone -q https://github.tools.sap/cs-devops/$name $name.tmp
  cd $name.tmp
  git checkout -q $tag
  cd ..
  rm -rf $name
  mkdir $name
  for s in $sources; do mv $name.tmp/$s $name; done
  rm -rf $name.tmp
}

name=gotpl
tag=$(go list -m -f '{{.Version}}' github.tools.sap/cs-devops/$name)
sources="cmd go.mod go.sum"
echo "Downloading github.tools.sap/cs-devops/$name $tag"
download "$name" "$tag" "$sources"

name=kubernetes-testing
tag=$(go list -m -f '{{.Version}}' github.tools.sap/cs-devops/$name)
sources="framework templates go.mod go.sum"
echo "Downloading github.tools.sap/cs-devops/$name $tag"
download "$name" "$tag" "$sources"

