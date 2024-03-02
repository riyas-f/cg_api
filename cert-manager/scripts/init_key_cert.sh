#!/bin/sh
set -e

mkdir cert

export COUNTRY="ID"
export STATE="West Java"
export LOCALITY="Banudng"
export ORGANIZATION="Aditya Software Ltd"
export ORG_UNIT="Aditya Software Certificate Manager"
export COMMON_NAME="cert-manager.aditya.com"
export ISSUER="/C=${COUNTRY}/ST=${STATE}/L=${LOCALITY}/O=${ORGANIZATION}/OU=${ORG_UNIT}/CN=${COMMON_NAME}"
# export ABSOLUTE_CERT_PATH="$(pwd)/${CERT_FILE_PATH}"

echo "$(openssl rand -base64 32)" | tr -d '\n' > /tmp/passphrase
echo "$(cat /tmp/passphrase)"

openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
openssl req -new -x509 -key private_key.pem -out $CERT_FILE_PATH -days 365 -subj "${ISSUER}"
openssl pkcs8 -topk8 -inform PEM -outform PEM -in private_key.pem -out $PRIVATE_KEY_PATH -passout file:/tmp/passphrase



rm private_key.pem