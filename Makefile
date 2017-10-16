
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

run:
	go build . && ./phts

.PHONY: frontend
frontend:
	cd frontend/admin && ng build -prod --base-href /admin/frontend -d /admin/frontend/ -dop -op ../../static

