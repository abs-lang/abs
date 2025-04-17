FROM golang:1.24

RUN apt-get update
RUN apt-get install bash make git curl jq nodejs npm -y
ENV CONTEXT=abs
COPY . /abs
WORKDIR /abs
RUN go install github.com/go-bindata/go-bindata/...@latest
RUN go mod vendor

CMD ["bash"]
