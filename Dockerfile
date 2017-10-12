FROM golang:1.8
ADD . /go/src/github.com/previousnext/m8s
WORKDIR /go/src/github.com/previousnext/m8s
RUN go get github.com/mitchellh/gox
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/previousnext/m8s/bin/m8s_linux_amd64 /usr/local/bin/m8s
CMD ["m8s"]
