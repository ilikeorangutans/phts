
frontend_admin_src_files := $(shell find frontend/admin/src -type f -name '*.ts' -or -name '*.html' -or -name '*.css')
frontend_admin_dist_dir := 'frontend/admin/dist/'
#frontend_admin_dist_files := $(addprefix $(frontend_admin_dist_dir),favicon.ico index.html inline.bundle.js inline.bundle.map main.bundle.js main.bundle.js.map polyfills.bundle.js polyfills.bundle.js.map styles.bundle.js styles.bundle.js.map vendor.bundle.js vendor.bundle.js.map)
frontend_admin_dist_files := frontend/admin/dist/main.bundle.js

.PHONY: foo
foo: $(frontend_admin_dist_files)
	@echo $(frontend_admin_dist_files)


$(frontend_admin_dist_files): $(frontend_admin_src_files)
	cd frontend/admin && ng build

.PHONY: test
test:
	go test . ./db ./model

.PHONY: frontend
frontend: frontend/admin/dist
	echo $(source_files)
	cp -rv frontend/admin/dist/* static/admin

frontend/admin/dist:
	cd frontend/admin && ng build

clean:
	rm -rf frontend/admin/dist

run:
	go build . && ./phts

