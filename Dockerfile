FROM golang:1.12

RUN apt-get update
RUN apt-get install bash make git curl jq -y
RUN go get github.com/c-bata/go-prompt

# The aim os to eventually remove these
RUN go get -v github.com/mattn/go-colorable
RUN go get -v github.com/mattn/go-tty

ENV CONTEXT=abs

COPY . /go/src/github.com/abs-lang/abs
WORKDIR /go/src/github.com/abs-lang/abs
RUN go get -d -v ./...

CMD bash
