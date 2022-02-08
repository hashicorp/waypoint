FROM golang:1.17.6

RUN apt-get update; apt-get install unzip

# Protoc
# TODO(izaak): discover the protoc version from the nix files
RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip -O /tmp/protoc.zip && \
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