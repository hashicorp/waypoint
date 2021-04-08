# syntax = docker.mirror.hashicorp.services/docker/dockerfile:experimental

#--------------------------------------------------------------------
# builder builds the Waypoint binaries
#--------------------------------------------------------------------

FROM docker.mirror.hashicorp.services/golang:alpine AS builder

RUN apk add --no-cache git gcc libc-dev openssh make

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
RUN go get github.com/kevinburke/go-bindata/...

COPY . /tmp/wp-src
WORKDIR /tmp/wp-src

RUN --mount=type=cache,target=/root/.cache/go-build make bin
RUN --mount=type=cache,target=/root/.cache/go-build make bin/entrypoint

#--------------------------------------------------------------------
# imgbase builds the "img" tool and all of its dependencies
#--------------------------------------------------------------------

# We build a fork of img for now so we can get the `img inspect` CLI
# Watch this PR: https://github.com/genuinetools/img/pull/324
FROM docker.mirror.hashicorp.services/golang:alpine AS imgbuilder

RUN apk add --no-cache \
	bash \
	build-base \
	gcc \
	git \
	libseccomp-dev \
	linux-headers \
	make

RUN git clone https://github.com/mitchellh/img.git /img
WORKDIR /img
RUN go get github.com/go-bindata/go-bindata/go-bindata
RUN make BUILDTAGS="seccomp noembed dfrunmount dfsecrets dfssh" && mv img /usr/bin/img

# Copied from img repo, see notes for specific reasons:
# https://github.com/genuinetools/img/blob/d858ac71f93cc5084edd2ba2d425b90234cf2ead/Dockerfile
FROM docker.mirror.hashicorp.services/alpine AS imgbase
RUN apk add --no-cache autoconf automake build-base byacc gettext gettext-dev \
    gcc git libcap-dev libtool libxslt img runc
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

# Notes on img and what is required to make it work, since there's a lot
# of small details below that are absolutely required for everything to
# come together:
#
#  - img, runc, newuidmap, newgidmap need to be installed
#  - libseccomp-dev must be installed for runc
#  - newuidmap/newgidmap need to have suid set (u+s)
#  - /etc/subuid and /etc/subgid need to have an entry for the user
#  - USER, HOME, and XDG_RUNTIME_DIR all need to be set
#

FROM docker.mirror.hashicorp.services/alpine

COPY --from=imgbuilder /usr/bin/img /usr/bin/img
COPY --from=imgbase /usr/bin/runc /usr/bin/runc
COPY --from=imgbase /usr/bin/newuidmap /usr/bin/newuidmap
COPY --from=imgbase /usr/bin/newgidmap /usr/bin/newgidmap

# libseccomp-dev is required for runc
# git is for gitrefpretty() and other calls for Waypoint
RUN apk add --no-cache libseccomp-dev git

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
RUN chmod u+s /usr/bin/newuidmap /usr/bin/newgidmap \
  && mkdir -p /run/user/100 \
  && chown -R waypoint /run/user/100 /home/waypoint \
  && echo waypoint:100000:65536 | tee /etc/subuid | tee /etc/subgid

USER waypoint
ENV USER waypoint
ENV HOME /home/waypoint
ENV XDG_RUNTIME_DIR=/run/user/100

ENTRYPOINT ["/usr/bin/waypoint"]
