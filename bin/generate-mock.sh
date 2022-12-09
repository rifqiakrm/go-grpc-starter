#!/usr/bin/env bash

set -euo pipefail

for file in `find . -name '*.go' | grep -v proto`; do
    if `grep -q 'interface {' ${file}`; then
        dest=${file//internal\//}
        mockgen -source=${file} -destination=test/mock/${dest}
    fi
done