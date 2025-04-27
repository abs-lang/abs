FROM golang:1.24

RUN apt-get update
RUN apt-get install bash make git curl jq nodejs npm -y
RUN cd /tmp && \
    wget https://github.com/upx/upx/releases/download/v5.0.0/upx-5.0.0-amd64_linux.tar.xz && \
    tar -xf upx-5.0.0-amd64_linux.tar.xz && \
    mv upx-5.0.0-amd64_linux/upx /usr/bin
ENV CONTEXT=abs
COPY . /abs
WORKDIR /abs
RUN go install github.com/go-bindata/go-bindata/...@latest
RUN go mod vendor

CMD ["bash"]
