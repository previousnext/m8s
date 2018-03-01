FROM golang:1.8
ADD . /go/src/github.com/previousnext/k8s-black-death
WORKDIR /go/src/github.com/previousnext/k8s-black-death
RUN go get github.com/mitchellh/gox
RUN make build

FROM alpine:3.6
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/previousnext/k8s-black-death/bin/k8s-black-death_linux_amd64 /usr/local/bin/k8s-black-death
CMD ["k8s-black-death"]
