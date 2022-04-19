# nested-logrus-formatter


all: test demo

prepare:
	@echo "Installing golangci-lint"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	@golangci-lint run ./...

dependency:
	@go get -v ./...

test: dependency
	@go test ./...

coverage: dependency
	go test . -v -covermode=count -coverprofile=coverage.out

demo:
	go run example/main.go

.PHONY: all prepare test coverage demo