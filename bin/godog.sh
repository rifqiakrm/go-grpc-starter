#!/usr/bin/env bash

set -euo pipefail

if [ -f ".env" ]; then
  cp .env.sample .env
else
  cp .env.sample .env
fi

for folder in `find . -name '*.feature' | sed -E 's|/[^/]+$||' | sort -u`; do
    godog run ${folder}/*.feature --format=cucumber > report.json
done