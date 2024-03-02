FROM golang:1.22-alpine AS build
ARG PORT
ARG PKI_HOST
RUN echo "port: ${PORT}"

RUN apk update && \
    apk add --no-cache curl

WORKDIR /github.com/AdityaP1502/Instant-Messanging/reverse_proxy
COPY . .
RUN go mod download
RUN curl -k https://$PKI_HOST/certificate -o /tmp/root-ca.crt
RUN --mount=type=cache,target=/root/.cache/go-build go build -o app main.go

FROM alpine:latest AS final
ARG SERVICE_NAME
ARG PORT

WORKDIR /app
COPY  --from=build /tmp/root-ca.crt /usr/local/share/ca-certificates/root-ca.crt

RUN cat /usr/local/share/ca-certificates/root-ca.crt >> /etc/ssl/certs/ca-certificates.crt 

COPY --from=build /github.com/AdityaP1502/Instant-Messanging/reverse_proxy/app .
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/reverse_proxy/config/app.config.json config/app.config.json

RUN apk --update --no-cache add curl

EXPOSE $PORT
CMD ["./app"]