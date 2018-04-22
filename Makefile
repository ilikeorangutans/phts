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
	go test ./test/integration/db ./test/integration/model

.PHONY: repl
repl:
	go run ./repl/main.go

run: phts
	DB_HOST=localhost DB_USER=phts DB_PASSWORD=secret DB_SSLMODE=false DB_NAME=phts ./phts

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
