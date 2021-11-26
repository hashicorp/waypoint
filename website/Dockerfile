FROM --platform=linux/amd64 docker.mirror.hashicorp.services/node:14.17.0-alpine 
# RUN apk add --update --no-cache git make g++ automake autoconf libtool nasm libpng-dev
RUN apk add --update --no-cache \
    autoconf \
    automake \
    bash \
    git \
    g++ \
    libtool \
    libc6-compat \
    libjpeg-turbo-dev \
    libpng-dev \
    make \
    nasm
WORKDIR /website
RUN npm install -g npm@latest
RUN mkdir /website/website-preview
COPY ./package.json /website/website-preview/package.json
COPY ./package-lock.json /website/website-preview/package-lock.json
RUN cd /website/website-preview && npm install --production=false