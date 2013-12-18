#!/bin/bash

WORKDIR=$(mktemp -d)

cleanup()
{
	rm -rf "$WORKDIR"
}

trap cleanup EXIT

pushd $WORKDIR >/dev/null
git clone /srv/git/website
cd website
# DONT_HIDE_FAILURES is a variable that is specific to our code. For easier
# development, the code does not require access to a postgresql database, but
# we want the build to fail in case there were postgresql issues.
if ! env -i LC_CTYPE=C.UTF-8 DONT_HIDE_FAILURES=1 PGUSER=anon /usr/local/bin/jekyll build
then
	echo "jekyll build failed"
	exit 1
fi

# All directories need to be group-writable, otherwise other users cannot
# create new files in them.
find _site -type d -exec chmod g+w '{}' \;
rsync --delete --delete-after --perms -v -r $WORKDIR/website/_site/ /var/www/www.noname-ev.de/_site/
popd >/dev/null
