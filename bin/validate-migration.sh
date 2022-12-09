#!/usr/bin/env bash

set -euo pipefail

for file in `find db -type f -name '*.sql'`; do
    if ! [[ ${file} =~ \/[0-9]{14}_ ]]; then
        echo "${file} doesn't use standard migration file name. the timestamp must exactly 14 characters long"
        echo "use command 'make migration module=[module-name] name=[migration-name]' to create a migration file"
        exit 1
    fi
done