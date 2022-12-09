#!/usr/bin/env bash

set -euo pipefail

make migrate-schema url=$1

for module in `ls db/migrations`; do
    make migrate module=${module} url=$1
done