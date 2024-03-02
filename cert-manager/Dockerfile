FROM golang:1.22-alpine AS build
ARG PORT
ARG CERT_FILE_PATH
ARG PRIVATE_KEY_PATH

RUN echo "port: ${PORT}"

RUN 

RUN apk update && \
    apk add --no-cache openssl

WORKDIR /github.com/AdityaP1502/Instant-Messanging/cert_manager
COPY . .
RUN scripts/init_key_cert.sh

RUN go mod download
RUN --mount=type=cache,target=/root/.cache/go-build go build -o app main.go

FROM alpine:latest AS final

ARG PORT
ARG CERT_FILE_PATH
ARG PRIVATE_KEY_PATH
ENV CERT_FILE_PATH ${CERT_FILE_PATH}
ENV PRIVATE_KEY_PATH ${PRIVATE_KEY_PATH}
ENV HOST "0.0.0.0:${PORT}"

WORKDIR /app


COPY --from=build /tmp/passphrase /tmp/passphrase
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/cert_manager/app .
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/cert_manager/cert/ cert/
EXPOSE $PORT
CMD ["./app"]