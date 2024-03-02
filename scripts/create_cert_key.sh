#!/bin/sh
set -e

mkdir -p service/${SERVICE_NAME}/cert 

cd service/${SERVICE_NAME}

export COUNTRY="ID"
export STATE="West Java"
export LOCALITY="Banudng"
export ORGANIZATION="Aditya Software Ltd"
export ORG_UNIT="Aditya Software Services"
export COMMON_NAME="${SERVICE_NAME}"
export UNIT="/C=${COUNTRY}/ST=${STATE}/L=${LOCALITY}/O=${ORGANIZATION}/OU=${ORG_UNIT}/CN=${COMMON_NAME}"
export SUBJECT_ALT_NAME="DNS.1:${SERVICE_NAME}"

echo $CSR_FILE_PATH
echo $CERT_FILE_PATH

echo "$(openssl rand -base64 32)" | tr -d '\n' > /tmp/passphrase
echo "$(cat /tmp/passphrase)"

# Generate private key
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
# openssl req -new -key private_key.pem -out $CSR_FILE_PATH -subj "${UNIT}" -config /tmp/extfile.cnf

openssl req -new \
    -key private_key.pem \
    -subj "${UNIT}" \
    -reqexts SAN \
    -config <(cat /etc/ssl/openssl.cnf \
        <(printf "\n[SAN]\nsubjectAltName=${SUBJECT_ALT_NAME}")) \
    -out $CSR_FILE_PATH

openssl req -noout -text -in $CSR_FILE_PATH

# send the certificate request to cert manager
curl -k https://$PKI_HOST/certificate -o cert/root-ca.crt

curl -k -v -F "csr=@${CSR_FILE_PATH}" https://$PKI_HOST/certificate/sign -o $CERT_FILE_PATH

openssl x509 -noout -text -in $CERT_FILE_PATH
 
# store key into pkcs8 format
openssl pkcs8 -topk8 -inform PEM -outform PEM -in private_key.pem -out $PRIVATE_KEY_PATH -passout file:/tmp/passphrase

rm $CSR_FILE_PATH
rm private_key.pem
