#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

DIR=$(dirname "${BASH_SOURCE[0]}")
BASEDIR="$DIR"/..
cd "$BASEDIR"

for f in $(find . -name "*.go"); do
  if [[ "$(sed -n 1p < $f)" != "/*" ]] || [[ "$(sed -n 2p < $f)" != Copyright* ]]; then 
    echo "Adding boilerplate: $f"
    cat hack/LICENSE_BOILERPLATE.txt > $f.tmp
    echo "" >> $f.tmp
    cat $f >> $f.tmp
    mv $f.tmp $f
  fi
done