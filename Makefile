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

.PHONY: frontend
frontend:
	cd frontend/admin && ng build -prod --base-href /admin/frontend -d /admin/frontend/ -dop -op ../../static

docker: dist
	docker build -t phts:latest -t phts:$(SHA) docker

.PHONY: dist-run
dist-run: dist
	DB_HOST=localhost DB_USER=phts DB_PASSWORD=secret DB_SSLMODE=false DB_NAME=phts ./phts

.PHONY: dist
dist: ui-dist phts-dist

ui-dist: ui
	mkdir -p docker/ui/dist
	cp ui/dist/* docker/ui/dist

phts-dist:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags $(DIST_LD_FLAGS) .
	mkdir -p docker/db/migrate
	cp phts docker
	cp db/migrate/* docker/db/migrate

.PHONY: ui
ui: ui/dist/index.html

UI_SOURCES=$(shell find ./ui/src/ -type f -iname '*.ts' -o -iname '*.css' -o -iname '*.html')

ui/dist/index.html: $(UI_SOURCES)
	cd ui  && ./node_modules/@angular/cli/bin/ng build -prod -aot -d static/

PHTS_SOURCES=$(shell find ./ -type f -iname '*.go')

phts: $(PHTS_SOURCES)
	go build .

.PHONY: clean
clean:
	rm -v phts
	rm -rf ui/dist
