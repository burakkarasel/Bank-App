#!/bin/sh

set -e

echo "run migrations"
/app/soda migrate

echo "start the app"
exec "$@"