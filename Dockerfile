FROM golang:1.10

RUN apt-get update
RUN apt-get install bash make git curl jq -y
RUN go get github.com/c-bata/go-prompt

# The aim os to eventually remove these
RUN go get -v github.com/mattn/go-colorable
RUN go get -v github.com/mattn/go-tty

COPY . /go/src/abs
WORKDIR /go/src/abs
RUN go get -d -v ./...

RUN chmod +x scripts/release.sh

CMD bash
