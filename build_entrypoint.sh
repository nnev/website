#!/bin/bash
set -e
cd /tmp/go/src/github.com/nnev/website/www
createdb nnev || true
./_createdb.sh
