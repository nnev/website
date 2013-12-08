#!/bin/sh

# Creates and fills a production database (on eris) with _schema.sql and _data.sql
sed 's/mero/postgres/g' _schema.sql | psql nnev
psql nnev < _data.sql

echo 'GRANT SELECT ON vortraege TO PUBLIC;' | psql nnev
echo 'GRANT SELECT ON termine TO PUBLIC;' | psql nnev
echo 'GRANT SELECT ON zusagen TO PUBLIC;' | psql nnev

echo 'GRANT ALL PRIVILEGES ON DATABASE nnev TO nnweb;' | psql nnev
echo 'GRANT ALL ON ALL SEQUENCES IN SCHEMA PUBLIC TO nnweb;' | psql nnev
