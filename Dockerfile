FROM golang:1.10-alpine

RUN apk add --update bash make
COPY . /go/src/abs
WORKDIR /go/src/abs
RUN go get -d -v ./...

CMD bash
