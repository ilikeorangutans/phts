
.PHONY: test

test:
	go test . ./db ./model

all-tests: test integration-test

integration-test:
	go test ./test/integration/db ./test/integration/model

.PHONY: repl
repl:
	go run ./repl/main.go

run:
	go build . && ./phts

