VERSION ?= "34.0"

build:
	go build -ldflags="-X 'weldr-client/cmd/composer-cli/root.Version=${VERSION}'" ./cmd/composer-cli

check:
	go vet ./... && golint -set_exit_status ./...

test:
	go test -v ./...

.PHONY: build check test
