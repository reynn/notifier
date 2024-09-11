RESULT_HASH := $(shell date -u -Iseconds)

build: build-proto build-server

build-server:
	$(MAKE) -C cmd/server build

build-proto:
	buf generate --clean --config proto/buf.yaml --template proto/buf.gen.yaml

run-server: build-server
	./notifier-server

test:
	@mkdir -p test-results
	go test -v -json ./... > test-results/TestResults-$(RESULT_HASH).json
