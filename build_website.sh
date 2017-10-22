#!/bin/bash
set -e

cd /tmp/go/src/github.com/nnev/website/www
DONT_HIDE_FAILURES=1 jekyll build
