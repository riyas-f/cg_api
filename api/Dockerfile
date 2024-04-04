FROM docker.io/library/golang:1.22-alpine@sha256:0466223b8544fb7d4ff04748acc4d75a608234bf4e79563bff208d2060c0dd79 AS build
ARG SERVICE_NAME
ARG PORT
ARG CSR_FILE_PATH
ARG CERT_FILE_PATH
ARG PRIVATE_KEY_PATH
ARG CSR_PATH
ARG PKI_HOST

ARG ROOT_CA_CRT="cert/root-ca.crt"

RUN echo "service_name: ${SERVICE_NAME}"
RUN echo "port: ${PORT}"

RUN apk update && \
    apk add --no-cache openssl curl

WORKDIR /github.com/AdityaP1502/Instant-Messanging/api/

COPY go.mod . 
COPY go.sum .
RUN go mod download

COPY scripts/create_cert_key.sh scripts/create_cert_key.sh
RUN chmod +x scripts/create_cert_key.sh
RUN scripts/create_cert_key.sh

COPY . .

# RUN chmod +x scripts/create_cert_key.sh
# RUN scripts/create_cert_key.sh

# RUN go mod download
RUN --mount=type=cache,target=/root/.cache/go-build go build -o app service/$SERVICE_NAME/main.go

FROM docker.io/library/alpine:3.19@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b AS final

ARG SERVICE_NAME
ARG PORT
ARG CERT_FILE_PATH
ARG PRIVATE_KEY_PATH

ENV CERT_FILE_PATH ${CERT_FILE_PATH}
ENV PRIVATE_KEY_PATH ${PRIVATE_KEY_PATH}
ENV ROOT_CA_CERT=/usr/local/share/ca-certificates/root-ca.crt

WORKDIR /app

COPY --from=build /github.com/AdityaP1502/Instant-Messanging/api/service/$SERVICE_NAME/cert/root-ca.crt /usr/local/share/ca-certificates/root-ca.crt
RUN cat /usr/local/share/ca-certificates/root-ca.crt >> /etc/ssl/certs/ca-certificates.crt 

# set env for root ca

COPY --from=build /tmp/passphrase /tmp/passphrase
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/api/scripts scripts
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/api/app .
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/api/service/$SERVICE_NAME/config/app.config.json config/app.config.json
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/api/service/$SERVICE_NAME/cert/ cert/

RUN chmod +x scripts/compose-entrypoint.sh
EXPOSE $PORT
CMD ["./app"]