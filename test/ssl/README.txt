README
=====

1. Generate wildcard self-signed certificate

1.1 Generate private key and csr
openssl req -newkey rsa:2048 -nodes -keyout go-docker-builder.test.key -out go-docker-builder.test.csr -config go-docker-builder.test.cnf

1.2 Generate certificate
openssl x509 -signkey go-docker-builder.test.key -in go-docker-builder.test.csr -req -days 365 -out go-docker-builder.test.crt
