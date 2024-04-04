FROM postgres:latest

ARG SERVICE_NAME

COPY scripts/ /scripts/
COPY sql/${SERVICE_NAME}-init.sql /docker-entrypoint-initdb.d/init.sql
