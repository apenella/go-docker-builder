#!/bin/sh
set -eu

# allow to authenticate and connect to gitserver
if [ -d "/root/.ssh" ]; then
    chmod 700 /root/.ssh
    ssh-keyscan -H gitserver >> /root/.ssh/known_hosts
    chmod 600 /root/.ssh/*
fi

# allow docker client to interact with registry.go-docker-builder.test
if [ -f /ssl/go-docker-builder.test.crt ]; then
    mkdir -p /etc/docker/certs.d/registry.go-docker-builder.test
    cp /ssl/go-docker-builder.test.crt /etc/docker/certs.d/registry.go-docker-builder.test/ca.crt
fi

/usr/local/bin/dockerd-entrypoint.sh 2> /dev/null &

while ! nc -z localhost 2376; do
    >&2 echo " Waiting for dockerd to be ready..."
    sleep 0.5 # wait for 1/2 of the second before check again
done

exec "$@"