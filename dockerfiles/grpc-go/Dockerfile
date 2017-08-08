FROM golang:1.8

RUN go get google.golang.org/grpc && \
    go get -u github.com/golang/protobuf/proto && \
    go get -u github.com/golang/protobuf/protoc-gen-go

ADD bin/protoc /usr/local/bin/protoc
RUN chmod a+x /usr/local/bin/protoc
