SHA=$(shell git rev-parse HEAD)
NOW=$(shell date +%FT%T%z)
DIST_LD_FLAGS="-X github.com/ilikeorangutans/phts/version.Sha=$(SHA) -X github.com/ilikeorangutans/phts/version.BuildTime=$(NOW)"

PHTS_SOURCES=$(shell find ./ -type f -iname '*.go')

.PHONY: test

test:
	go test . ./db ./model

install:
	go install ./...

all-tests: test integration-test

integration-test:
	go test ./test/integration/db ./test/integration/model

setup-integration-test-env:
	docker run -it --rm -e PGPASSWORD=secret --network env_default --link env_db_1:postgres postgres psql -h postgres -U \
		phts -c "CREATE DATABASE phts_test"
	docker run -it --rm -e PGPASSWORD=secret --network env_default --link env_db_1:postgres postgres psql -h postgres -U \
		phts -c "CREATE ROLE phts_test WITH LOGIN PASSWORD 'phts'; GRANT ALL PRIVILEGES ON DATABASE phts_test TO phts_test;"

################################################################################
# dev environment
################################################################################

DEV_DB_NAME=phts_dev
DEV_DB_USER=phts_dev
DEV_DB_PASSWORD=secret
MINIO_ACCESS_KEY=minio
MINIO_SECRET_KEY=supersecret
MINIO_BUCKET=phts-dev

.PHONY: start-env stop-env setup-dev-env run wipe-dev-env start-psql

start-env: stop-env
	docker run --rm --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=$(DEV_DB_PASSWORD) -d postgres
	docker run --rm --name minio -d -e MINIO_ACCESS_KEY=$(MINIO_ACCESS_KEY) -e MINIO_SECRET_KEY=$(MINIO_SECRET_KEY) -p 9000:9000 minio/minio server /data
	echo "waiting for minio to come up..."
	sleep 5
	mcli config host add localhost http://localhost:9000/ $(MINIO_ACCESS_KEY) $(MINIO_SECRET_KEY)

stop-env:
	-docker stop postgres
	-docker stop minio

setup-dev-env: wipe-dev-env
	-docker exec -i -t postgres psql -U postgres -c "create role $(DEV_DB_USER) with login password '$(DEV_DB_PASSWORD)';"
	-docker exec -i -t postgres psql -U postgres -c "create database $(DEV_DB_NAME) with owner $(DEV_DB_USER);"
	-mcli mb localhost/$(MINIO_BUCKET)

wipe-dev-env:
	-mcli rb localhost/$(MINIO_BUCKET)
	-docker exec -i -t postgres psql -U postgres -c "drop database $(DEV_DB_NAME);"
	-rm -r tmp

start-psql:
	docker exec -i -t postgres psql -U postgres $(DEV_DB_NAME)

.PHONY: repl
repl:
	go run ./repl/main.go

run: phts
	PHTS_SERVER_URL=http://localhost:8080 PHTS_DB_HOST=localhost PHTS_ADMIN_EMAIL=test@test.local PHTS_ADMIN_PASSWORD=test PHTS_DB_USER=$(DEV_DB_USER) PHTS_DB_PASSWORD=$(DEV_DB_PASSWORD) PHTS_DB_SSLMODE=false PHTS_DB_DATABASE=$(DEV_DB_NAME) PHTS_STORAGE_ENGINE=minio PHTS_MINIO_ENDPOINT=localhost:9000 PHTS_MINIO_ACCESS_KEY=$(MINIO_ACCESS_KEY) PHTS_MINIO_SECRET_KEY=$(MINIO_SECRET_KEY) PHTS_MINIO_BUCKET=$(MINIO_BUCKET) ./phts

phts: $(PHTS_SOURCES)
	go build -ldflags $(DIST_LD_FLAGS) -tags debug ./cmd/phts/

################################################################################
# dist targets
################################################################################

.PHONY: dist-all
dist-all: target/linux-amd64/phts docker-linux-arm

.PHONY: ui-dist
ui-dist: admin-ui-dist
	$(MAKE) -C ui dist

.PHONY: admin-ui-dist
admin-ui-dist:
	$(MAKE) -C ui-admin dist

target/%/: ui-dist
	mkdir -p $(@)
	cp -r ui-admin/dist $(@)/ui-admin
	mkdir -p $(@)/ui
	cp -r ui/dist $(@)/ui
	# TODO there's a bug here if this reruns it'll copy static into static
	cp -r static $(@)/static
	mkdir -p $(@)/db/migrate
	mkdir -p $(@)/templates/services/internal
	cp db/migrate/* $(@)/db/migrate
	cp -vr templates/services/internal/* $(@)/templates/services/internal
	cp docker/Dockerfile $(@)/

target/linux-amd64/phts: target/linux-amd64/ $(PHTS_SOURCES)
	mkdir -p target/linux-amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags $(DIST_LD_FLAGS) -o target/linux-amd64/phts ./cmd/phts

target/linux-arm/phts: target/linux-arm/ $(PHTS_SOURCES)
	mkdir -p target/linux-arm
	# see https://github.com/golang/go/wiki/GoArm
	GOARM=7 GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -ldflags $(DIST_LD_FLAGS) -o target/linux-arm/phts ./cmd/phts

.PHONY: setup-buildx
setup-buildx:
	docker run --rm --privileged docker/binfmt:820fdd95a9972a5308930a2bdfb8573dd4447ad3
	docker buildx rm arm-builder
	docker buildx create --name arm-builder
	docker buildx inspect --bootstrap arm-builder


################################################################################
# docker targets
################################################################################

DOCKER_TAGS=phts:latest phts:$(SHA) registry.ilikeorangutans.me/apps/phts:latest registry.ilikeorangutans.me/apps/phts:$(SHA)

.PHONY: docker-linux-arm
docker-linux-arm: target/linux-arm/phts
	#docker buildx build $(DOCKER_TAGS) --platform linux/arm/v7,linux/amd64 -f docker/Dockerfile target/linux-arm --load
	# TODO docker buildx can't build multi arch and load them at the time :|
	docker buildx build $(foreach tag,$(DOCKER_TAGS),-t $(tag) ) --platform linux/arm/v7 -f docker/Dockerfile target/linux-arm --load

################################################################################
# clean targets
################################################################################

.PHONY: ui-clean
ui-clean:
	$(MAKE) -C ui-admin clean
	$(MAKE) -C ui clean

.PHONY: clean
clean:
	-rm -rv phts docker/ui-admin docker/phts docker/db target
