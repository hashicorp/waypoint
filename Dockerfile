# syntax = hashicorp.jfrog.io/docker/docker/dockerfile:experimental

#--------------------------------------------------------------------
# builder builds the Waypoint binaries
#--------------------------------------------------------------------

FROM hashicorp.jfrog.io/docker/golang:alpine AS builder

RUN apk add --no-cache git gcc libc-dev openssh

RUN mkdir -p /tmp/wp-prime
COPY go.sum /tmp/wp-prime
COPY go.mod /tmp/wp-prime

WORKDIR /tmp/wp-prime

RUN mkdir -p -m 0600 ~/.ssh \
    && ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
RUN --mount=type=ssh --mount=type=secret,id=ssh.config --mount=type=secret,id=ssh.key \
    GIT_SSH_COMMAND="ssh -o \"ControlMaster auto\" -F \"/run/secrets/ssh.config\"" \
    go mod download

COPY . /tmp/wp-src

WORKDIR /tmp/wp-src

RUN apk add --no-cache make
RUN go get github.com/kevinburke/go-bindata/...
RUN --mount=type=cache,target=/root/.cache/go-build make bin
RUN --mount=type=cache,target=/root/.cache/go-build make bin/entrypoint

#--------------------------------------------------------------------
# imgbase builds the "img" tool and all of its dependencies
#--------------------------------------------------------------------

# Copied from img repo, see notes for specific reasons:
# https://github.com/genuinetools/img/blob/d858ac71f93cc5084edd2ba2d425b90234cf2ead/Dockerfile
FROM hashicorp.jfrog.io/docker/alpine AS imgbase
RUN apk add --no-cache autoconf automake build-base byacc gettext gettext-dev \
    gcc git libcap-dev libtool libxslt img
RUN git clone https://github.com/shadow-maint/shadow.git /shadow
WORKDIR /shadow
RUN git checkout 59c2dabb264ef7b3137f5edb52c0b31d5af0cf76
RUN ./autogen.sh --disable-nls --disable-man --without-audit \
    --without-selinux --without-acl --without-attr --without-tcb \
    --without-nscd \
    && make \
    && cp src/newuidmap src/newgidmap /usr/bin

#--------------------------------------------------------------------
# final image
#--------------------------------------------------------------------

FROM hashicorp.jfrog.io/docker/alpine

COPY --from=imgbase /usr/bin/img /usr/bin/img
COPY --from=imgbase /usr/bin/newuidmap /usr/bin/newuidmap
COPY --from=imgbase /usr/bin/newgidmap /usr/bin/newgidmap

COPY --from=builder /tmp/wp-src/waypoint /usr/bin/waypoint
COPY --from=builder /tmp/wp-src/waypoint-entrypoint /usr/bin/waypoint-entrypoint

VOLUME ["/data"]

RUN addgroup waypoint && \
    adduser -S -G waypoint waypoint && \
    mkdir /data/ && \
    chown -R waypoint:waypoint /data

USER waypoint

ENTRYPOINT ["/usr/bin/waypoint"]
