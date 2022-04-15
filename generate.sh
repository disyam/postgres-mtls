#!/bin/bash

#  create certs directory if not exists
mkdir -p certs

# generate root private key
openssl ecparam -name prime256v1 -genkey -noout -out certs/ca.key

# create root certificate
openssl req -x509 -key certs/ca.key -subj "/CN=root" -days 3650 -out certs/ca.crt

# generate server private key
openssl ecparam -name prime256v1 -genkey -noout -out certs/server.key

# create server CSR
openssl req -new -key certs/server.key -subj "/CN=server" -addext "subjectAltName=DNS:127.0.0.1" -out certs/server.csr

# create server certificate
openssl x509 -req -in certs/server.csr -days 3650 -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -copy_extensions copy -out certs/server.crt

# generate client private key
openssl ecparam -name prime256v1 -genkey -noout -out certs/client.key

# create client CSR
openssl req -new -key certs/client.key -subj "/CN=postgres" -out certs/client.csr

# create client certificate
openssl x509 -req -in certs/client.csr -days 3650 -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/client.crt
