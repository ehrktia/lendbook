#! /usr/bin/bash
clear
set -euo pipefail
file_count=`ls -lart | grep *.go | wc -l`
if (($file_count>0)); then
echo "removing mock files from ${ROOT_DIR}/mocks"
rm ${ROOT_DIR}/mocks/*.go
fi
echo  "generating mocks in ${ROOT_DIR}/mocks..."
mockery --dir internal --outpkg mocks --output mocks --all \
    --exclude internal/graph --with-expecter=true
echo "mock files successfully generated."
ls -lart ${ROOT_DIR}/mocks
