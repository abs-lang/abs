FROM golang:1.15

RUN apt-get update
RUN apt-get install bash make git curl jq -y
ENV CONTEXT=abs
COPY . /abs
WORKDIR /abs
RUN go mod vendor

CMD bash
