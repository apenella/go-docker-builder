#!/bin/sh
set -eu

if [ -d "/root/.ssh" ]; then
    chmod 700 /root/.ssh
    ssh-keyscan -H gitserver >> /root/.ssh/known_hosts
    chmod 600 /root/.ssh/*
fi

/usr/local/bin/dockerd-entrypoint.sh &

while ! nc -z localhost 2376; do
    echo " Waiting for dockerd to be ready..."
    sleep 0.5 # wait for 1/10 of the second before check again
done

exec "$@"