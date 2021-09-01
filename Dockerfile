FROM golang:1.17

RUN apt-get update
RUN apt-get install bash make git curl jq nodejs npm -y
RUN go get -u github.com/jteeuwen/go-bindata/...
ENV CONTEXT=abs
COPY . /abs
WORKDIR /abs
RUN go mod vendor

CMD bash
