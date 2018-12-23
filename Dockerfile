FROM golang:1.11

RUN apt-get update
RUN apt-get install bash make git curl jq -y
RUN go get github.com/c-bata/go-prompt

# The aim os to eventually remove these
RUN go get -v github.com/mattn/go-colorable
RUN go get -v github.com/mattn/go-tty

COPY . /go/src/github.com/abs-lang/abs
WORKDIR /go/src/github.com/abs-lang/abs
RUN go get -d -v ./...

RUN chmod +x scripts/release.sh
# This is simply done because Go
# will build faster. A docker build
# is probably less frequent than an ABS
# build, so...
RUN ./scripts/release.sh

CMD bash
