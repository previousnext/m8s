#!/usr/bin/make -f

export CGO_ENABLED=0

PROJECT=github.com/previousnext/skpr
MOCKERNETES_HOST=127.0.0.1
MOCKERNETES_PORT=8080

default: build

# Builds the project.
build: generate
	gox -os='linux darwin' -arch='amd64' -output='bin/skpr_{{.OS}}_{{.Arch}}' -ldflags='-extldflags "-static"' $(PROJECT)

# Run all lint checking with exit codes for CI
lint: generate
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting
test: generate
	nc -z ${MOCKERNETES_HOST} ${MOCKERNETES_PORT}
	go test -cover ./...

# Generates code.
generate:
	go generate

VERSION=$(shell git describe --tags --always)

# Releases the project Docker Hub.
release-github: build
	ghr -u previousnext "${VERSION}" ./bin/

release: build release-github

GITBOOK=gitbook
DOCS_SRC=./docs
DOCS_DST=./docs/_book
docs-build:
	$(GITBOOK) install ${DOCS_SRC}
	$(GITBOOK) build --log=debug ${DOCS_SRC} ${DOCS_DST}

.PHONY: default build lint test generate release release-github docs-build
