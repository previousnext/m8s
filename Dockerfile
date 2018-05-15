FROM golang:1.8
RUN go get github.com/mitchellh/gox
ADD . /go/src/github.com/previousnext/m8s
WORKDIR /go/src/github.com/previousnext/m8s
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/previousnext/m8s/bin/m8s_linux_amd64 /usr/local/bin/m8s

ENTRYPOINT ["m8s"]
CMD ["server"]
