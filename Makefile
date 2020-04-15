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

start-db:
	cd env && docker-compose up -d

stop-db:
	cd env && docker-compose down

start-psql:
	docker run -it --rm -e PGPASSWORD=secret --network env_default --link env_db_1:postgres postgres psql -h postgres -U phts

setup-integration-test-env:
	docker run -it --rm -e PGPASSWORD=secret --network env_default --link env_db_1:postgres postgres psql -h postgres -U \
		phts -c "CREATE DATABASE phts_test"
	docker run -it --rm -e PGPASSWORD=secret --network env_default --link env_db_1:postgres postgres psql -h postgres -U \
		phts -c "CREATE ROLE phts_test WITH LOGIN PASSWORD 'phts'; GRANT ALL PRIVILEGES ON DATABASE phts_test TO phts_test;"

DEV_DB_NAME=phts_dev
DEV_DB_USER=phts_dev
DEV_DB_PASSWORD=secret

setup-dev-env:
	docker run -it --rm -e PGPASSWORD=secret --network env_default --link env_db_1:postgres postgres psql -h postgres -U \
		phts -c "CREATE DATABASE $(DEV_DB_NAME)"
	docker run -it --rm -e PGPASSWORD=secret --network env_default --link env_db_1:postgres postgres psql -h postgres -U \
		phts -c "CREATE ROLE $(DEV_DB_USER) WITH LOGIN PASSWORD '$(DEV_DB_PASSWORD)'; GRANT ALL PRIVILEGES ON DATABASE $(DEV_DB_NAME) TO $(DEV_DB_USER);"

.PHONY: repl
repl:
	go run ./repl/main.go

run: phts
	DB_HOST=localhost DB_USER=$(DEV_DB_USER) DB_PASSWORD=$(DEV_DB_PASSWORD) DB_SSLMODE=false DB_NAME=$(DEV_DB_NAME) ./phts

phts: $(PHTS_SOURCES)
	go build .

################################################################################
# dist targets
################################################################################

.PHONY: dist-all
dist-all: target/linux-amd64/phts docker-linux-arm

.PHONY: ui-dist
ui-dist: admin-ui-dist public-ui-dist

.PHONY: admin-ui-dist
admin-ui-dist:
	$(MAKE) -C ui-admin dist

.PHONY: public-ui-dist
public-ui-dist:
	$(MAKE) -C ui-public dist

target/%/: ui-dist
	mkdir -p $(@)
	cp -r ui-admin/dist $(@)/ui-admin
	cp -r ui-public/dist $(@)/ui-public
	mkdir -p $(@)/db/migrate
	cp db/migrate/* $(@)/db/migrate
	cp docker/Dockerfile $(@)/

target/linux-amd64/phts: target/linux-amd64/ $(PHTS_SOURCES)
	mkdir -p target/linux-amd64
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags $(DIST_LD_FLAGS) -o target/linux-amd64/phts .

target/linux-arm/phts: target/linux-arm/ $(PHTS_SOURCES)
	mkdir -p target/linux-arm
	# see https://github.com/golang/go/wiki/GoArm
	GOARM=7 GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -ldflags $(DIST_LD_FLAGS) -o target/linux-arm/phts .

################################################################################
# docker targets
################################################################################

DOCKER_TAGS=phts:latest phts:$(SHA) registry.ilikeorangutans.me/apps/phts:latest registry.ilikeorangutans.me/apps/phts:$(SHA)

.PHONY: docker-arm
docker-linux-arm: target/linux-arm/phts
	#docker buildx build $(DOCKER_TAGS) --platform linux/arm/v7,linux/amd64 -f docker/Dockerfile target/linux-arm --load
	# TODO docker buildx can't build multi arch and load them at the time :|
	docker buildx build $(foreach tag,$(DOCKER_TAGS),-t $(tag) ) --platform linux/arm/v7 -f docker/Dockerfile target/linux-arm --load

################################################################################
# clean targets
################################################################################

.PHONY: ui-clean
ui-clean:
	$(MAKE) -C ui-public clean
	$(MAKE) -C ui-admin clean

.PHONY: clean
clean:
	-rm -rv phts docker/ui-admin docker/ui-public docker/phts docker/db target
