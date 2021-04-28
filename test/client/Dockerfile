# Based on https://github.com/timberio/docker-client/blob/master/Dockerfile
ARG GOLANG_VERSION="1.15"
ARG DOCKER_VERSION="20.10"

FROM golang:${GOLANG_VERSION}-alpine AS golang

FROM docker:${DOCKER_VERSION}-dind

COPY --from=golang /usr/local/go /usr/local/go


ENV GOPATH /go
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
ENV PATH /usr/local/go/bin:$PATH
ENV PATH $GOPATH/bin:$PATH 
WORKDIR $GOPATH

COPY entrypoint.sh /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["go","version"]
