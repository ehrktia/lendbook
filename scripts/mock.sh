#! /usr/bin/bash
clear
set -euo pipefail
echo -n "generating mocks..."
mockery --dir internal --outpkg mocks --output mocks --all --with-expecter=true
