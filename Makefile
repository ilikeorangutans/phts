
.PHONY: test

test:
	go test . ./db ./model

run:
	go build . && ./phts

