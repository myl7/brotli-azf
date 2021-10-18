#!/bin/bash
set -euo pipefail

export CC=musl-gcc
go build -o bin/brotli --ldflags '-linkmode external -extldflags "-static"'
