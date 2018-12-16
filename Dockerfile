FROM golang:1.10-alpine

RUN apk add --update bash make git
RUN go get github.com/c-bata/go-prompt
COPY . /go/src/abs
WORKDIR /go/src/abs
RUN go get -d -v ./...

CMD bash
