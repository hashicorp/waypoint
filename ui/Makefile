ROOT:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

all: build

deps: node_modules clean

clean:
	rm -rf ./tmp

build: deps
	yarn run build

start: deps
	yarn run start

test: deps test-node
	yarn run test

test-ember: deps test-node
	yarn run test:ember

lint: deps
	yarn run lint:hbs && yarn run lint:js

node_modules: yarn.lock package.json
	yarn install

.PHONY: all deps build start test test-ember lint clean