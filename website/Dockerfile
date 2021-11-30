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
# note: /website/package.json is copied only for sake
# of consistency in use of start script
COPY package.json /website/package.json
ADD scripts /website/scripts
ADD website-preview /website/website-preview
WORKDIR /website
RUN npm install -g npm@latest
# install dependencies within /website/website-preview
RUN cd /website/website-preview && npm install --production=false