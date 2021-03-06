SHARED_SOURCES=$(shell find ./projects/shared -type f)
ADMIN_SOURCES=$(shell find ./projects/admin -type f)
FRONTED_SOURCES=$(shell find ./projects/frontend -type f)
SERVICES_INTERNAL_SOURCES=$(shell find ./projects/services-internal -type f)
PROJECT_SOURCES=angular.json Makefile package.json tsconfig.json tslint.json yarn.lock

NG_BINARY=./node_modules/.bin/ng

.PHONY: run-admin run run-services-internal run-frontend dist

# --------------------------------------------------------------------------------
# dev targets

.PHONY: lint
lint:
	$(NG_BINARY) lint

# --------------------------------------------------------------------------------
# run targets

run-admin:
	$(NG_BINARY) serve admin --watch --aot --base-href /admin

run-services-internal:
	$(NG_BINARY) serve services-internal --watch --aot --base-href /services/internal

run-frontend:
	$(NG_BINARY) serve frontend --watch --aot --base-href /

# --------------------------------------------------------------------------------
# dist targets

.PHONY: yarn
yarn:
	yarnpkg

dist: ./dist/admin ./dist/services-internal ./dist/frontend

dist/admin: yarn $(PROJECT_SOURCES) $(ADMIN_SOURCES) $(SHARED_SOURCES)
	$(NG_BINARY) build admin --aot --build-optimizer --optimization --prod --deploy-url /admin/static/ --base-href /admin

dist/services-internal: yarn $(PROJECT_SOURCES) $(SERVICES_INTERNAL_SOURCES) $(SHARED_SOURCES)
	$(NG_BINARY) build services-internal --aot --build-optimizer --optimization --prod --deploy-url /services/internal/static/ --base-href /services/internal

dist/frontend: yarn $(PROJECT_SOURCES) $(FRONTEND_SOURCES) $(SHARED_SOURCES)
	$(NG_BINARY) build frontend --aot --build-optimizer --optimization --prod --deploy-url /static/ --base-href /

# --------------------------------------------------------------------------------
# clean targets

.PHONY: clean

clean:
	-rm -r dist/*
