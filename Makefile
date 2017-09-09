
.PHONY: test

test:
	go test . ./db ./model

integration-test:
	go test -v ./test/integration/db

run:
	go build . && ./phts

