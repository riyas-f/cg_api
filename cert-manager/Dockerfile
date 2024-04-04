FROM docker.io/library/golang:1.22-alpine@sha256:0466223b8544fb7d4ff04748acc4d75a608234bf4e79563bff208d2060c0dd79 AS build
ARG PORT
ARG CERT_FILE_PATH
ARG PRIVATE_KEY_PATH 
ARG GCP_CERT_FILE_BUCKET_URL
ARG GCP_PRIVATE_KEY_SECRET_NAME
ARG GCP_PRIVATE_KEY_PASSPHRASE_SECRET_NAM 

RUN echo "port: ${PORT}"

RUN 

RUN apk update && \
    apk add --no-cache curl openssl

WORKDIR /github.com/AdityaP1502/Instant-Messanging/cert_manager
COPY . .

# RUN chmod +x scripts/init_key_cert.sh
# RUN scripts/init_key_cert.sh

RUN go mod download
RUN --mount=type=cache,target=/root/.cache/go-build go build -o app main.go

FROM docker.io/library/alpine:3.19@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b AS final

ARG PORT
# ARG CERT_FILE_PATH
# ARG PRIVATE_KEY_PATH
# ENV CERT_FILE_PATH ${CERT_FILE_PATH}
# ENV PRIVATE_KEY_PATH ${PRIVATE_KEY_PATH}
ENV HOST "0.0.0.0:${PORT}"

WORKDIR /app


COPY --from=build /tmp/passphrase /tmp/passphrase
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/cert_manager/app .
# COPY --from=build /github.com/AdityaP1502/Instant-Messanging/cert_manager/cert/ cert/
EXPOSE $PORT
CMD ["./app"]