FROM alpine:3.7

RUN mkdir -p /app/ui/dist
RUN mkdir -p /app/tmp
WORKDIR /app

COPY phts .
COPY ui/dist/* ui/dist/

VOLUME /app/tmp

ARG DB_HOST
ENTRYPOINT ["/app/phts"]

