#!/usr/bin/make -f

export CGO_ENABLED=0

PROJECT=github.com/previousnext/m8s
VERSION=$(shell git describe --tags --always)
COMMIT=$(shell git rev-list -1 HEAD)
BUILD=$(shell date)

# Builds the project
build: generate
	gox -os='linux darwin' \
	    -arch='amd64' \
	    -output='bin/m8s_{{.OS}}_{{.Arch}}' \
	    -ldflags='-extldflags "-static" -X github.com/previousnext/m8s/cmd.GitVersion=${VERSION} -X github.com/previousnext/m8s/cmd.GitCommit=${COMMIT}' \
	    $(PROJECT)

# Generate any necessary code.
generate:
	go generate

# Run all lint checking with exit codes for CI
lint: generate
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting
test: generate
	go test -cover ./...

IMAGE=previousnext/m8s

# Releases the project Docker Hub
release-docker:
	# Building M8s...
	docker build -t ${IMAGE}:${VERSION} -t ${IMAGE}:latest .
	docker push ${IMAGE}:${VERSION}
	docker push ${IMAGE}:latest
	# Building M8s UI...
	docker build -t ${IMAGE}-ui:${VERSION} -t ${IMAGE}-ui:latest ui
	docker push ${IMAGE}-ui:${VERSION}
	docker push ${IMAGE}-ui:latest

release-github: build
	ghr -u previousnext "${VERSION}" ./bin/

release: release-docker release-github