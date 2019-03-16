#!/usr/bin/make -f

export CGO_ENABLED=0

PROJECT=github.com/previousnext/m8s

# Builds the project
build: generate
	gox -os='linux darwin' -arch='amd64' -output='bin/m8s_{{.OS}}_{{.Arch}}' -ldflags='-extldflags "-static"' $(PROJECT)

# Generate any necessary code.
generate:
	go generate

# Run all lint checking with exit codes for CI
lint: generate
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting
test: generate
	go test -cover ./server/...
	go test -cover ./cmd/...

run: build
	bin/m8s_linux_amd64 server --port=8443 \
	                           --token=123456789 \
			           --kubeconfig=.kube/config

IMAGE=previousnext/m8s
VERSION=$(shell git describe --tags --always)

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

PROTOBUF=$(PWD)/pb

# Generates a new Protobuf Golang package
protobuf:
	rm -fR $(PROTOBUF)
	mkdir -p $(PROTOBUF)
	docker run -it -w $(PWD) -v $(PWD):$(PWD) nickschuch/grpc-go:latest /bin/bash -c 'protoc -I . m8s.proto --go_out=plugins=grpc:$(PROTOBUF)'

.PHONY: build lint test release-docker release-github release protobuf
