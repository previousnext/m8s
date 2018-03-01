#!/usr/bin/make -f

export CGO_ENABLED=0

PROJECT=github.com/previousnext/gopher
VERSION=$(shell git describe --tags --always)
COMMIT=$(shell git rev-list -1 HEAD)

# Builds the project.
build:
	gox -os='linux darwin' \
	    -arch='amd64' \
	    -output='bin/gopher_{{.OS}}_{{.Arch}}' \
	    -ldflags='-extldflags "-static" -X github.com/previousnext/gopher/cmd.GitVersion=${VERSION} -X github.com/previousnext/gopher/cmd.GitCommit=${COMMIT}' \
	    $(PROJECT)

# Run all lint checking with exit codes for CI.
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting.
test:
	go test -cover ./...

IMAGE=previousnext/gopher

# Releases the project Docker Hub.
release-docker:
	docker build -t ${IMAGE}:${VERSION} -t ${IMAGE}:latest .
	docker push ${IMAGE}:${VERSION}
	docker push ${IMAGE}:latest

release-github: build
	ghr -u previousnext "${VERSION}" ./bin/

release: release-docker release-github

.PHONY: build lint test release-docker release-github release
