ARG GOLANG_VERSION="1.15"

FROM golang:${GOLANG_VERSION}-alpine AS golang

WORKDIR /
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o go-dummy-app
