FROM golang:1.19

ARG PROTOC_VERSION="3.17.3"

RUN apt-get update; apt-get install unzip

# Protoc
# TODO(izaak): discover the protoc version from the nix files
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-$(uname -m | sed s/aarch64/aarch_64/g).zip -O /tmp/protoc.zip && \
    unzip /tmp/protoc.zip -d /tmp && \
    mv /tmp/bin/protoc /usr/local/bin/ && \
    chmod +x /usr/local/bin/protoc && \
    mv /tmp/include/* /usr/local/include/

# Copy files required to update tooling
RUN mkdir -p /tools/tools
COPY ./tools/tools.go /tools/tools
COPY ./Makefile /tools
COPY go.mod /tools
COPY go.sum /tools

RUN make -C /tools tools

WORKDIR /waypoint
