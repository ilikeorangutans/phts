SHA=$(shell git rev-parse HEAD)
NOW=$(shell date +%FT%T%z)
DIST_LD_FLAGS="-X github.com/ilikeorangutans/phts/version.Sha=$(SHA) -X github.com/ilikeorangutans/phts/version.BuildTime=$(NOW)"

.PHONY: test

test:
	go test . ./db ./model

install:
	go install ./...

all-tests: test integration-test

integration-test:
	PHTS_MIGRATION_DIR=$(shell pwd)/db/migrate go test ./test/integration/db ./test/integration/model ./services/admin

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

.PHONY: docker
docker: dist
	docker build -t phts:latest -t phts:$(SHA) docker

.PHONY: dist
dist: admin-ui-dist public-ui-dist phts-dist

.PHONY: admin-ui-dist
admin-ui-dist: admin-ui
	mkdir -p docker/ui-admin
	cp -r ui-admin/dist docker/ui-admin

.PHONY: public-ui-dist
public-ui-dist: public-ui
	mkdir -p docker/ui-public
	cp -r ui-public/dist docker/ui-public

.PHONY: admin-ui
admin-ui:
	$(MAKE) -C ui-admin dist

.PHONY: public-ui
public-ui:
	$(MAKE) -C ui-public dist

phts-dist:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags $(DIST_LD_FLAGS) .
	mkdir -p docker/db/migrate
	cp phts docker
	cp db/migrate/* docker/db/migrate

PHTS_SOURCES=$(shell find ./ -type f -iname '*.go')

phts: $(PHTS_SOURCES)
	go build .

.PHONY: ui-clean
ui-clean:
	$(MAKE) -C ui-public clean
	$(MAKE) -C ui-admin clean

.PHONY: clean
clean:
	rm -rv phts docker/ui-admin docker/ui-public docker/phts docker/db
