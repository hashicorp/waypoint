# Setting up Kind and Kubernetes

## What is this?

This folder contains a helper script to help users get started with setting up
`kind` and `kubernetes` locally. Because Waypoint requires a kubernetes (k8s) cluster
with a load balancer configured, and because setting this up is a bit non-trivial,
this folder exists to ease the pain of bringing up a kubernetes cluster locally
so users can get started quickly without figuring all of this out on their own.

## The Plan

- Use and set up https://kind.sigs.k8s.io/ locally
  + Follow the quick start guide to get more familiar with kind
  + https://kind.sigs.k8s.io/docs/user/quick-start/

## Requirements

- Docker installed locally
- Golang 1.11+ installed locally
- kubectl installed and available on the path
  + https://kubernetes.io/docs/tasks/tools/install-kubectl/
- A way to automate and setup a realistic but minimalist environment to test Waypoint with

### macOS and Windows Steps

This is much easier for these platforms. Docker Desktop provides a simple way
to setup and run kubernetes. This is the recommended approach.

### Linux Steps

#### Automated

Run the script inside this folder to automatically setup k8s with kind and metallb.

It will eventually ask you to write up the IP Address range for metallb based
on the networked container created. Follow the instructions to set the address,
and the rest should be taken care of.

```bash
./setup-k8s.sh
```

After this script runs, you should be ready to run a `waypoint install` for
the kubernetes platform!

#### Manual

_this section is a work in progress_

1) kind create cluster --config configs/cluster-config.yaml
2) kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.9.3/manifests/namespace.yaml
3) kubectl create secret generic -n metallb-system memberlist --from-literal=secretkey="$(openssl rand -base64 128)"
4) kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.9.3/manifests/metallb.yaml
5) Get docker subnet from networked container `docker ps -a`, then `docker inspect <container_id>`, and update metallb addresses in `configs/metallb-config.yaml` to represent your local docker subnet
6) kubectl apply -f configs/metallb-config.yaml

### Optional steps??

7) clone demo app `codyde/hashi-demo-app`
8) kubectl apply -f namespace.yaml
9) kubectl apply -f kubernetes-demoapp.yaml

### Setup waypoint

Now you are ready to install the waypoint server to your local kind k8s cluster

### Debugging k8s

Just some useful `kubectl` commands for determining what's going on with your
local k8s cluster.

```
kubectl get svc -A
```

```
kubectl get all
```

Inspect a deployed application in a pod

```
kubectl describe pod/example-nodejs-01eqxfhphddst35xb04pp4m2gs-6f559cb4bd-gcfp5  
```
