#!/bin/bash
# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: BUSL-1.1


# NOTE(briancain): This script uses 'kind' to automatically bring up a local
# kubernetes cluster. It includes optional support for configuring an
# Ingress controller node.

set -o errexit

setupIngress=$WP_K8S_INGRESS

# NOTE(briancain): This version is both the kind node version as well as the
# version of Kubernetes that kind sets up. In this case, we'd be installing
# kubernetes version 1.22.
kindVersion="${KIND_NODE_VERSION:-kindest/node:v1.22.7}"

echo "Setting up local docker registry..."
echo

reg_name='kind-registry'
reg_port='5000'
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --name "${reg_name}" \
    registry:2 | 2>/dev/null
fi

echo "Setting up kubernetes with kind and metallb..."
echo

if [ -n "${setupIngress}" ]; then
  echo "Creating kind cluster with ingress controller..."
  echo
cat <<EOF | kind create cluster --image="$kindVersion" --config=-
  kind: Cluster
  apiVersion: kind.x-k8s.io/v1alpha4
  containerdConfigPatches:
  - |-
    [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
        endpoint = ["http://${reg_name}:${reg_port}"]
  nodes:
  - role: control-plane
    kubeadmConfigPatches:
    - |
      kind: InitConfiguration
      nodeRegistration:
        kubeletExtraArgs:
          node-labels: "ingress-ready=true"
    extraPortMappings:
    - containerPort: 80
      hostPort: 80
      listenAddress: "0.0.0.0"
      protocol: TCP
    - containerPort: 443
      hostPort: 443
      listenAddress: "0.0.0.0"
      protocol: TCP
EOF

else
  echo "Creating kind cluster with cluster-config.yaml..."
  echo
cat <<EOF | kind create cluster --image="$kindVersion" --config=-
  kind: Cluster
  apiVersion: kind.x-k8s.io/v1alpha4
  containerdConfigPatches:
  - |-
    [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
        endpoint = ["http://${reg_name}:${reg_port}"]
EOF

fi

echo "Connecting registry to cluster network..."
echo
# connect the registry to the cluster network
# (the network may already be connected)
docker network connect "kind" "${reg_name}" || true
# Document the local registry
# https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${reg_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF

echo
echo "Applying metallb namespace..."
echo
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.11.0/manifests/namespace.yaml

echo "Create secret for metallb-system node..."
echo
kubectl create secret generic -n metallb-system memberlist --from-literal=secretkey="$(openssl rand -base64 128)"

echo "Applying metallb manifest..."
echo
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.11.0/manifests/metallb.yaml

echo
echo
echo $"Obtaining container IP to set as range in metallb-config.yaml..."
CONTAINERID=$(docker ps -a --filter="expose=6443" -q)
IPADDR=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $CONTAINERID)

echo "IP Address of networked container is ${IPADDR}"
echo

echo "Automatically picking an IP range based on networked contianer IP..."

IPADDR_RANGE="$(awk -F"." '{print $1"."$2"."$3".20"}'<<<$IPADDR)-$(awk -F"." '{print $1"."$2"."$3".50"}'<<<$IPADDR)"

#read -p "Enter a range to define for metallb IP Addresses based on this container (like 172.18.0.20-172.18.0.50):`echo $'\n> '` " IPADDR_RANGE
echo "IP Address range is ${IPADDR_RANGE}"
echo

sed s/%ADDR_RANGE%/$IPADDR_RANGE/g \
   configs/metallb-config-template.yaml > configs/metallb-config-set.yaml

echo "Applying metallb-config-set.yaml with ip address range applied..."
kubectl apply -f configs/metallb-config-set.yaml

if [ -n "${setupIngress}" ]; then
  echo
  echo "Setting up NGINX ingress controller..."
  echo

  kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

  # NOTE: uncomment this if you wish to use contour
  #kubectl apply -f https://projectcontour.io/quickstart/contour.yaml
  #kubectl patch daemonsets -n projectcontour envoy -p '{"spec":{"template":{"spec":{"nodeSelector":{"ingress-ready":"true"},"tolerations":[{"key":"node-role.kubernetes.io/master","operator":"Equal","effect":"NoSchedule"}]}}}}'
fi

echo
echo "Done! You should be ready to 'waypoint install -platform=kubernetes -accept-tos' on a local kubernetes!"
exit 0
