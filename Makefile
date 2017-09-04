
.PHONY: test

test:
	go test . ./db ./model

.PHONY: frontend
frontend: frontend/admin/dist
	cp -rv frontend/admin/dist/* static/admin

frontend/admin/dist:
	cd frontend/admin && ng build --prod

clean:
	rm -rf frontend/admin/dist

run:
	go build . && ./phts

