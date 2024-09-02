build-server:
	$(MAKE) -C cmd/server build

build-proto:
	buf generate --clean --config proto/buf.yaml --template proto/buf.gen.yaml


run-server: build-server
	./notifier-server
