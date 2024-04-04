# add sha256 hash into the image to sped up load image
FROM docker.io/library/node:16-alpine@sha256:a1f9d027912b58a7c75be7716c97cfbc6d3099f3a97ed84aa490be9dee20e787


ARG PORT
ARG SERVICE_NAME
ARG CSR_FILE_PATH
ARG CERT_FILE_PATH
ARG PRIVATE_KEY_PATH
ARG CSR_PATH
ARG PKI_HOST

RUN apk update && \
    apk add --no-cache openssl curl

# install node js dependencies 
# independent of codechanges
# thanks to: https://bitjudo.com/blog/2014/03/13/building-efficient-dockerfiles-node-dot-js/

COPY package.json /tmp/package.json
RUN cd /tmp && npm install
RUN mkdir -p /app && cp -a /tmp/node_modules app/

WORKDIR /app

# create cert key
COPY scripts/create_cert_key.sh scripts/create_cert_key.sh
RUN chmod +x scripts/create_cert_key.sh
RUN chomod +x scripts/compose-entrypoint.sh
RUN scripts/create_cert_key.sh

# copy the application code into
# the directory
COPY . .

ENV port ${PORT}
ENV CERT_FILE_PATH ${CERT_FILE_PATH}
ENV PRIVATE_KEY_PATH ${PRIVATE_KEY_PATH}
ENV ROOT_CA_CERT=/usr/local/share/ca-certificates/root-ca.crt
EXPOSE $PORT

CMD ["node", "app.js"]