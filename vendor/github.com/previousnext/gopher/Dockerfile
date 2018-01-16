FROM previousnext/golang:1.9
ADD . /go/src/github.com/previousnext/gopher
WORKDIR /go/src/github.com/previousnext/gopher
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/previousnext/gopher/bin/gopher_linux_amd64 /usr/local/bin/gopher
CMD ["gopher"]
