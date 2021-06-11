FROM base-registry.go-docker-builder.test:5000/alpine:3.13

RUN apk update \
    && apk upgrade

RUN dd if=/dev/urandom of=/f1 bs=250MB count=1
RUN dd if=/dev/urandom of=/f2 bs=80MB count=1
RUN dd if=/dev/urandom of=/f3 bs=90MB count=1
RUN dd if=/dev/urandom of=/f4 bs=100MB count=1
RUN dd if=/dev/urandom of=/f5 bs=120MB count=1