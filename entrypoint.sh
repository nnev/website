#!/bin/bash
set -e
echo "Fetching website and build dependencies"
go get -u github.com/nnev/website/...
echo "Updating and installing website from /usr/src"
cp -r /usr/src/* /tmp/go/src/github.com/nnev/website/
go install github.com/nnev/website/...
cd /tmp/go/src/github.com/nnev/website/www
createdb nnev || true
./_createdb.sh
termine -hook=/build_website.sh -connect="dbname=nnev host=$PGHOST sslmode=disable" next 4
jekyll build
exec /usr/bin/supervisord -n -c /etc/supervisor/supervisord.conf
