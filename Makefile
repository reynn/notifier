RESULT_HASH := $(shell date -u -Iseconds)
GIT_SHA := $(shell git rev-parse HEAD)
DOCKER_REGISTRY := ghcr.io/reynn/notifier

build: build-proto build-server build-client

build-server:
	$(MAKE) -C cmd/server build

build-client:
	$(MAKE) -C cmd/client build

build-proto:
	buf generate --clean --config proto/buf.yaml --template proto/buf.gen.yaml

docker-build:
	docker build --build-arg GIT_SHA=$(GIT_SHA) -t $(DOCKER_REGISTRY)/notifier-server .

run-server: build-server
	./notifier-server

test:
	@mkdir -p test-results
	go test -v -json ./... > test-results/TestResults-$(RESULT_HASH).json
