#!/usr/bin/make -f

export CGO_ENABLED=0

PROJECT=github.com/previousnext/m8s

# Builds the project
build:
	gox -os='linux' -arch='amd64' -output='bin/m8s_{{.OS}}_{{.Arch}}' -ldflags='-extldflags "-static"' $(PROJECT)

# Run all lint checking with exit codes for CI
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting
test:
	go test -cover ./server/...
	go test -cover ./cmd/...

IMAGE=previousnext/m8s
VERSION=$(shell git describe --tags --always)

# Releases the project Docker Hub
release:
	docker build -t ${IMAGE}:${VERSION} .
	docker push ${IMAGE}:${VERSION}

PROTOBUF=$(PWD)/pb

# Generates a new Protobuf Golang package
protobuf:
	rm -fR $(PROTOBUF)
	mkdir -p $(PROTOBUF)
	docker run -it -w $(PWD) -v $(PWD):$(PWD) nickschuch/grpc-go:latest /bin/bash -c 'protoc -I . m8s.proto --go_out=plugins=grpc:$(PROTOBUF)'

.PHONY: build lint test release protobuf
