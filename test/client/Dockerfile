# Based on https://github.com/timberio/docker-client/blob/master/Dockerfile
ARG GOLANG_VERSION="1.15"

FROM golang:${GOLANG_VERSION}-alpine

ARG DOCKER_CLI_VERSION="20.10.5"
ENV DOWNLOAD_URL="https://download.docker.com/linux/static/stable/x86_64/docker-$DOCKER_CLI_VERSION.tgz"

RUN apk --no-cache add \
        curl \
        build-base \
    && curl -L $DOWNLOAD_URL | tar -xz docker \
    && mv docker/docker /usr/local/bin/docker

WORKDIR /app

ENTRYPOINT ["tail"]
CMD ["-f","/dev/null"]
