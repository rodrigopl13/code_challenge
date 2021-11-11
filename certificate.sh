#!/bin/sh

password=tmp_password

openssl genrsa -des3 -passout pass:$password -out localhost.pass.key 2048
openssl rsa -passin pass:$password -in localhost.pass.key -out localhost.key
rm localhost.pass.key
openssl req -new -key localhost.key -out localhost.csr \
    -subj "/C=MX/ST=Jalisco/L=Guadalajara/O=jobsity/OU=challenge/CN=localhost"
openssl x509 -req -days 10 -in localhost.csr -signkey localhost.key -out localhost.crt
