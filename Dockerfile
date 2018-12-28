FROM golang:1.11

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

# This is simply done because Go
# will build faster. A docker build
# is probably less frequent than an ABS
# build, so...
RUN make build_simple
RUN ./builds/abs ./scripts/release.abs

CMD bash
