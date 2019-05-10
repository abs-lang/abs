FROM golang:1.12

RUN apt-get update
RUN apt-get install bash make git curl jq -y
ENV CONTEXT=abs
COPY . /abs
WORKDIR /abs
RUN go mod vendor

CMD bash
