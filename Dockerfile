FROM golang:1.8-alpine

RUN apk add --no-cache git

RUN mkdir -p /go/src/github.com/aep/vault-init
WORKDIR /go/src/github.com/aep/vault-init
COPY *.go /go/src/github.com/aep/vault-init/
RUN go get -v
RUN go build

CMD ["./vault-init"]
