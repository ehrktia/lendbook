#! /usr/bin/bash
clear
set -euo pipefail
echo -n "remove old files"
rm mocks/*.go
echo -n "generating mocks..."
mockery --dir internal --outpkg mocks --output mocks --all \
    --exclude internal/graph --with-expecter=true
