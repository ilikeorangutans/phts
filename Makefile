
.PHONY: test

test:
	go test . ./db ./model

integration-test:
	go test -v ./test/integration

run:
	go build . && ./phts

