#!/usr/bin/make -f

VERSION=$(shell git describe --tags --always)

release:
	docker build -f dockerfiles/api/Dockerfile -t previousnext/pr-api:${VERSION} .
	docker build -f dockerfiles/cli/Dockerfile -t previousnext/pr-cli:${VERSION} .

GRPC_GO_IMAGE="nickschuch/skipper-grpc-go:latest"
GRPC_GO_TARGET="workspace/src/github.com/previousnext/pr/pb"
GRPC_RUN=docker run -it -w /data -v $(PWD):/data

protobuf:
	docker build -f dockerfiles/grpc-go/Dockerfile -t $(GRPC_GO_IMAGE) dockerfiles/grpc-go
	rm -fR $(GRPC_GO_TARGET) && mkdir -p $(GRPC_GO_TARGET)
	$(GRPC_RUN) $(GRPC_GO_IMAGE) /bin/bash -c 'protoc -I . pr.proto --go_out=plugins=grpc:$(GRPC_GO_TARGET)'

.PHONY: build push protobuf
