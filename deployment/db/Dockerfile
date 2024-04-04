FROM postgres:latest

ARG SERVICE_NAME
COPY sql/${SERVICE_NAME}-init.sql /docker-entrypoint-initdb.d/init.sql
COPY scripts /scripts
