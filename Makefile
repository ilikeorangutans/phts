
.PHONY: test

test:
	go test . ./db ./model

integration-test:
	go test -v ./integration_test

run:
	go build . && ./phts

