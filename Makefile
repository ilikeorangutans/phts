
.PHONY: test

test:
	go test . ./db ./model

all-tests: test integration-test

integration-test:
	go test -v ./test/integration/db ./test/integration/model

run:
	go build . && ./phts

