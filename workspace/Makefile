#!/usr/bin/make -f

# This allows us to have a GOPATH per peroject.
# Ensuring that we are all building with the same dependencies.
GOPATH := $(shell pwd):$(shell pwd)/vendor
export GOPATH

NAME=pr
PACKAGE=github.com/previousnext/$(NAME)

# Build binaries for linux/amd64 and darwin/amd64
build:
	gox -os='linux darwin' -arch='amd64' -output='bin/$(NAME)-api_{{.OS}}_{{.Arch}}' -ldflags='-extldflags "-static"' $(PACKAGE)/api
	gox -os='linux darwin' -arch='amd64' -output='bin/$(NAME)-cli_{{.OS}}_{{.Arch}}' -ldflags='-extldflags "-static"' $(PACKAGE)/cli

# Run all lint checking with exit codes for CI
lint:
	golint -set_exit_status $(PACKAGE)/...

# Run tests with coverage reporting
test:
	go test -cover $(PACKAGE)/...
