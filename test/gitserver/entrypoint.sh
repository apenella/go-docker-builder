#!/bin/sh

keys=${GIT_KEYS:-/git/keys}
git_ssh_folder="${GIT_SSH_FOLDER:-/home/git/.ssh}"

find "${keys}" -type f -name '*.pub' | while read key
do
    echo "Loading key: $key"
    cat "${key}" >> "${git_ssh_folder}/authorized_keys"
done

/usr/sbin/sshd -D
