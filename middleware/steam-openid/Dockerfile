FROM node:16-alpine

ARG SERVICE_NAME
ARG CSR_FILE_PATH
ARG CERT_FILE_PATH
ARG PRIVATE_KEY_PATH
ARG CSR_PATH
ARG PKI_HOST

RUN apk update && \
    apk add --no-cache openssl curl
    
WORKDIR /app

COPY . .
RUN scripts/create_cert_key.sh

COPY . . 

RUN npm install

ENV CERT_FILE_PATH ${CERT_FILE_PATH}
ENV PRIVATE_KEY_PATH ${PRIVATE_KEY_PATH}
ENV ROOT_CA_CERT=/usr/local/share/ca-certificates/root-ca.crt

CMD ["node", "app.js"]