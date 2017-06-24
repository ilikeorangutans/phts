
.PHONY:
test:
	go test . ./db ./model

.PHONY:
run:
	go build . && ./phts

