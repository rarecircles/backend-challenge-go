#!/bin/sh

set -e

go build -o build/challenge cmd/challenge/main.go cmd/challenge/logging.go 

echo "Build app successfully"