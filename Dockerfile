FROM golang:alpine

ADD . /go/src/github.com/johnweldon/block_landing

RUN go install github.com/johnweldon/block_landing

WORKDIR /go/src/github.com/johnweldon/block_landing

ENTRYPOINT /go/bin/block_landing

EXPOSE 9000
