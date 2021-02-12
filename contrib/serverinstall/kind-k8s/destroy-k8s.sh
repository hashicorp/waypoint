#!/bin/bash

kind delete cluster
docker rm kind-registry -f
