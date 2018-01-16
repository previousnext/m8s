FROM previousnext/golang:1.8

RUN go get -u \
        github.com/golang/dep/cmd/dep \
        github.com/golang/lint/golint

VOLUME /go/src/github.com/previousnext/gopher
WORKDIR /go/src/github.com/previousnext/gopher
