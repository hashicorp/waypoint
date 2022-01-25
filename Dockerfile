# syntax = docker.mirror.hashicorp.services/docker/dockerfile:experimental

#--------------------------------------------------------------------
# builder builds the Waypoint binaries
#--------------------------------------------------------------------

FROM docker.mirror.hashicorp.services/golang:1.17.5-alpine3.15 AS builder

RUN apk add --no-cache git gcc libc-dev make

RUN mkdir -p /tmp/wp-prime
COPY go.sum /tmp/wp-prime
COPY go.mod /tmp/wp-prime

WORKDIR /tmp/wp-prime

RUN go mod download
RUN go install github.com/kevinburke/go-bindata/go-bindata

COPY . /tmp/wp-src
WORKDIR /tmp/wp-src

RUN --mount=type=cache,target=/root/.cache/go-build make bin
RUN --mount=type=cache,target=/root/.cache/go-build make bin/entrypoint

# This is only used by ODR
FROM docker.mirror.hashicorp.services/busybox:stable-musl as busybox
RUN touch /tmp/.keep

#--------------------------------------------------------------------
# odr image
#--------------------------------------------------------------------
# This target is explicitly invoked from the command line, it's not used
# by the non-odr stages.
FROM gcr.io/kaniko-project/executor:v1.6.0 as odr

COPY --from=builder /tmp/wp-src/waypoint /kaniko/waypoint
COPY --from=busybox /bin/busybox /kaniko/busybox
COPY --from=busybox /tmp /kaniko/tmp

# We add busybox and populate it with the tool links to make the image
# easier to use (having a shell, basic tools, etc)
RUN ["/kaniko/busybox", "mkdir", "/kaniko/bin"]
RUN ["/kaniko/busybox", "--install", "-s", "/kaniko/bin"]

# Need to add the dir with our tools in PATH
ENV PATH $PATH:/kaniko/bin
ENV TMPDIR /kaniko/tmp

ENTRYPOINT ["/kaniko/waypoint"]

#--------------------------------------------------------------------
# final image
#--------------------------------------------------------------------

FROM docker.mirror.hashicorp.services/alpine:3.15.0

# git is for gitrefpretty() and other calls for Waypoint
RUN apk add --no-cache git

COPY --from=builder /tmp/wp-src/waypoint /usr/bin/waypoint
COPY --from=builder /tmp/wp-src/waypoint-entrypoint /usr/bin/waypoint-entrypoint

VOLUME ["/data"]

# NOTE: userid must be 100 here. Otherwise upgrades will fail due to user not
# having the proper permissions to read the server db due to a different userid
RUN addgroup waypoint && \
    adduser -S -u 100 -G waypoint waypoint && \
    mkdir /data/ && \
    chown -R waypoint:waypoint /data

# configure newuidmap/newgidmap to work with our waypoint user
RUN mkdir -p /run/user/100 \
  && chown -R waypoint /run/user/100 /home/waypoint \
  && echo waypoint:100000:65536 | tee /etc/subuid | tee /etc/subgid

USER waypoint
ENV USER waypoint
ENV HOME /home/waypoint
ENV XDG_RUNTIME_DIR=/run/user/100

ENTRYPOINT ["/usr/bin/waypoint"]
