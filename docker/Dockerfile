FROM alpine:latest

RUN apk update && apk add ca-certificates

RUN mkdir -p /app
RUN mkdir -p /app/tmp
RUN mkdir -p /app/db/migrate
RUN mkdir -p /app/templates/services/internal
WORKDIR /app

# TODO properly copy new ui files over
COPY phts .
COPY ui-admin/dist/* ui-admin/dist/
COPY ui/dist/ ui/dist/
COPY db/migrate/* db/migrate/
COPY templates/services/internal/ templates/services/internal
COPY static static

VOLUME /app/tmp

ARG DB_HOST
ENTRYPOINT ["/app/phts"]

