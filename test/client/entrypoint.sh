#!/bin/sh
set -eu

/usr/local/bin/dockerd-entrypoint.sh &

exec "$@"