FROM golang:1.22-alpine AS build
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
COPY . .

RUN scripts/create_cert_key.sh

RUN go mod download
RUN --mount=type=cache,target=/root/.cache/go-build go build -o app service/$SERVICE_NAME/main.go

FROM alpine:latest AS final

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
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/api/app .
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/api/service/$SERVICE_NAME/config/app.config.json config/app.config.json
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/api/service/$SERVICE_NAME/cert/ cert/

EXPOSE $PORT
CMD ["./app"]