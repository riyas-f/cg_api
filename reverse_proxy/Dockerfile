FROM golang:1.22-alpine AS build
ARG PORT
ARG PKI_HOST
ARG CERT_FILE_PATH
ARG PRIVATE_KEY_PATH
ARG CSR_FILE_PATH
RUN echo "port: ${PORT}"

RUN apk update && \
    apk add --no-cache openssl curl

WORKDIR /github.com/AdityaP1502/Instant-Messanging/reverse_proxy

COPY go.mod . 
COPY go.sum .
RUN go mod download

COPY scripts/create_cert_key.sh scripts/create_cert_key.sh
RUN scripts/create_cert_key.sh

COPY . .

# RUN curl -k https://$PKI_HOST/certificate -o /tmp/root-ca.crt
RUN --mount=type=cache,target=/root/.cache/go-build go build -o app main.go

FROM alpine:latest AS final
ARG SERVICE_NAME
ARG PORT
ARG CERT_FILE_PATH
ARG PRIVATE_KEY_PATH

ENV CERT_FILE_PATH ${CERT_FILE_PATH}
ENV PRIVATE_KEY_PATH ${PRIVATE_KEY_PATH}
ENV ROOT_CA_CERT=/usr/local/share/ca-certificates/root-ca.crt

WORKDIR /app
COPY  --from=build /tmp/root-ca.crt /usr/local/share/ca-certificates/root-ca.crt

RUN cat /usr/local/share/ca-certificates/root-ca.crt >> /etc/ssl/certs/ca-certificates.crt 

COPY --from=build /tmp/passphrase /tmp/passphrase
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/reverse_proxy/app .
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/reverse_proxy/config/app.config.json config/app.config.json
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/reverse_proxy/cert/ cert/

RUN apk --update --no-cache add curl

EXPOSE $PORT
CMD ["./app"]