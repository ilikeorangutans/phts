FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /app
RUN mkdir -p /app/tmp
RUN mkdir -p /app/db/migrate
WORKDIR /app

COPY phts .
COPY ui-admin/* ui-admin/
COPY ui-public/* ui-public/
COPY db/migrate/* db/migrate/

VOLUME /app/tmp

ARG DB_HOST
ENTRYPOINT ["/app/phts"]

