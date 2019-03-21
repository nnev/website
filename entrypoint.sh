#!/bin/bash
set -e
cd /tmp/go/src/github.com/nnev/website/www
termine -hook=/build_website.sh -connect="dbname=nnev host=$PGHOST sslmode=disable" next 4
jekyll build
exec /usr/bin/supervisord -n -c /etc/supervisor/supervisord.conf
