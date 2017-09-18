#!/usr/bin/make -f

VERSION=$(shell git describe --tags --always)
IMAGE=previousnext/m8s

release:
	docker build -f dockerfiles/api/Dockerfile -t ${IMAGE}:${VERSION} .
	docker push ${IMAGE}:${VERSION}

GRPC_GO_IMAGE="nickschuch/skipper-grpc-go:latest"
GRPC_GO_TARGET="workspace/src/github.com/previousnext/m8s/pb"
GRPC_RUN=docker run -it -w /data -v $(PWD):/data

protobuf:
	docker build -f dockerfiles/grpc-go/Dockerfile -t $(GRPC_GO_IMAGE) dockerfiles/grpc-go
	rm -fR $(GRPC_GO_TARGET) && mkdir -p $(GRPC_GO_TARGET)
	$(GRPC_RUN) $(GRPC_GO_IMAGE) /bin/bash -c 'protoc -I . m8s.proto --go_out=plugins=grpc:$(GRPC_GO_TARGET)'

.PHONY: build push protobuf
