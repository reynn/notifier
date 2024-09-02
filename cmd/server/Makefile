SERVICE_NAME := notifier
MODULE := server
GIT_SHA := $(shell git rev-parse HEAD)

build:
	go build \
		-o ../../$(SERVICE_NAME)-$(MODULE)\
		-ldflags="-X 'github.com/reynn/notifier/internal/constants.AppVersion=$(GIT_SHA)'" \
		main.go
