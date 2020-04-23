# phts

## What is it?

A photo sharing server with a focus on privacy that is easy to self host on even the smallest machines like a raspberry pi.

## What's it built on?

The backend app is written in Go and connects to postgresql databases. Different storage backends are supported: filesystem, google cloud storage, and minio (s3 compatible). The frontend is written in Angular.

## How can I run it?

TODO: add docker-compose files to easily run this.

## Env Variables

phts is configured primarily through environment variables. They're all prefixed with `PHTS_`.

- **PHTS_ADMIN_EMAIL** email address used as the user name for services/internal
- **PHTS_ADMIN_PASSWORD** password for services/internal
- **PHTS_BIND** bind address, defaults to `:8080`
- **PHTS_DB_SSL** use ssl to connect to the database
- **PHTS_DB_HOST** postgresql database host
- **PHTS_DB_USER** database username
- **PHTS_DB_PASSWORD** database password
- **PHTS_DB_DATABASE** database to connect to
- **PHTS_STORAGE_ENGINE** storage backend to use, supported values are `file`, `minio`, and `gcs`
- **PHTS_MINIO_ENDPOINT** minio endpoint with port
- **PHTS_MINIO_ACCESS_KEY** minio access key
- **PHTS_MINIO_SECRET_KEY** minio secret key
- **PHTS_MINIO_BUCKET** minio bucket

## Development

### Requirements

You'll need:

* [git](https://git-scm.com/)
* [Docker](https://www.docker.com/products/docker-desktop)
* [Docker buildx](https://docs.docker.com/buildx/working-with-buildx/) (if you're creating release builds)
* [make](https://www.gnu.org/software/make/)
* [go](https://golang.org/)
* [minio cli](https://github.com/minio/mc)
* [yarn](https://yarnpkg.com/) (for the UI)

### Starting the Dev Environment

You can run a complete phts installation locally. You'll need a postgres database and a minio server, both of which can automatically be set up by running:

```
$ make start-env
```
This make target will start postgres and minio in a docker container. Now you're ready to to set up the databases and minio bucket:
```
$ make setup-dev-env
```
If you need to wipe the environment clean you can do so by rerunning `make setup-dev-env` which will wipe the database and delete everything in the storage bucket.

If you need to inspect the database itself you can start a psql shell with `make start-psql`.

### Starting the Dev Server

Once your environment is up and running, you can start the actual phts server with
```
$ make run
```

This will build the backend app and start it, binding to `http://localhost:8080`.

## Release Builds

Needs docker and buildx for cross compiling. We could probably do without but if I ever want to build arm specific code in the container it'll make things easier. To set up a builder:
```
docker run --rm --privileged docker/binfmt:820fdd95a9972a5308930a2bdfb8573dd4447ad3
docker buildx rm arm-builder
docker buildx create --name arm-builder
docker buildx inspect --bootstrap arm-builder
```

To build release binaries and docker images:

```
make dist-all -j4
```

## Stuff to look at

Icons from https://fontawesome.com/
Spinners from https://github.com/tobiasahlin/SpinKit
