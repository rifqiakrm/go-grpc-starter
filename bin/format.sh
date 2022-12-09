#!/usr/bin/env bash

set -euo pipefail

for f in `find . -name '*.go'`; do
  # Defensive, just in case.
  if [[ -f ${f} ]]; then
    awk '/^import \($/,/^\)$/{if($0=="")next}{print}' ${f} > /tmp/file
    mv /tmp/file ${f}
  fi
done

goimports -w -local grpc-starter $(go list -f {{.Dir}} ./...)
gofmt -s -w .
