#!/bin/sh

# Creates and fills a production database (on eris) with _schema.sql and _data.sql
sed 's/mero/postgres/g' _schema.sql | psql nnev
psql nnev < _data.sql

echo 'GRANT SELECT ON vortraege TO PUBLIC;' | psql nnev
echo 'GRANT SELECT ON termine TO PUBLIC;' | psql nnev
echo 'GRANT SELECT ON zusagen TO PUBLIC;' | psql nnev

# for yarpnarp(1)
echo 'GRANT UPDATE ON zusagen TO nnweb;' | psql nnev
echo 'GRANT INSERT ON zusagen TO nnweb;' | psql nnev
echo 'GRANT DELETE ON zusagen TO nnweb;' | psql nnev

# for c14h(1)
echo 'GRANT UPDATE ON vortraege TO nnweb;' | psql nnev
echo 'GRANT INSERT ON vortraege TO nnweb;' | psql nnev
echo 'GRANT DELETE ON vortraege TO nnweb;' | psql nnev
