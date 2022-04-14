FROM golang:1.18

RUN apt-get update
RUN apt-get install bash make git curl jq nodejs npm -y
ENV CONTEXT=abs
COPY . /abs
WORKDIR /abs
RUN go install github.com/jteeuwen/go-bindata/...
RUN go mod vendor

CMD bash
