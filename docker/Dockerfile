FROM alpine:3.7

RUN mkdir -p /app/ui/dist
RUN mkdir -p /app/tmp
RUN mkdir -p /app/db/migrate
WORKDIR /app

COPY phts .
COPY ui/dist/* ui/dist/
COPY db/migrate/* db/migrate/

VOLUME /app/tmp

ARG DB_HOST
ENTRYPOINT ["/app/phts"]

