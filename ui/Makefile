SHARED_SOURCES=$(shell find projects/shared -type f)
ADMIN_SOURCES=$(shell find projects/admin -type f)
FRONTED_SOURCES=$(shell find projects/frontend -type f)
SERVICES_INTERNAL_SOURCES=$(shell find projects/services-internal -type f)
PROJECT_SOURCES=angular.json Makefile package.json tsconfig.json tslint.json yarn.lock

NG_BINARY=./node_modules/.bin/ng

.PHONY: run-admin run run-services-internal run-frontend dist

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

dist: dist/admin dist/services-internal dist/frontend

dist/admin: $(PROJECT_SOURCES) $(ADMIN_SOURCES) $(SHARED_SOURCES)
	$(NG_BINARY) build admin --aot --build-optimizer --optimization --prod --base-href /admin

dist/services-internal: $(PROJECT_SOURCES) $(SERVICES_INTERNAL_SOURCES) $(SHARED_SOURCES)
	$(NG_BINARY) build services-internal --aot --build-optimizer --optimization --prod --base-href /services/internal

dist/frontend: $(PROJECT_SOURCES) $(FRONTEND_SOURCES) $(SHARED_SOURCES)
	$(NG_BINARY) build frontend --aot --build-optimizer --optimization --prod --base-href /

# --------------------------------------------------------------------------------
# clean targets

.PHONY: clean

clean:
	-rm -r dist/*

