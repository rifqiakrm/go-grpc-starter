#!/usr/bin/env bash

set -euo pipefail

if [ ! -z "`git status -s | grep -e '.go'`" ]; then
  echo "Import blocks are not beautifully formatted for these files:"
  git status -s | grep -e '.go'
  echo "Run 'make pretty' or 'make format' before commit and push"
  exit 1
fi
