#!/bin/bash

echo "Setting up kubernetes with kind and metallb..."
echo

echo "Creating kind cluster with cluster-config.yaml..."
kind create cluster --config configs/cluster-config.yaml

echo "Applying metallb namespace..."
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.9.3/manifests/namespace.yaml

echo "Create secret for metallb-system node..."
kubectl create secret generic -n metallb-system memberlist --from-literal=secretkey="$(openssl rand -base64 128)"

echo "Applying metallb manifest..."
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.9.3/manifests/metallb.yaml

echo
echo
echo $"Obtaining container IP to set as range in metallb-config.yaml..."
CONTAINERID=$(docker ps -a --filter="expose=6443" -q)
IPADDR=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $CONTAINERID)

echo "IP Address of networked container is ${IPADDR}"
echo
read -p "Enter a range to define for metallb IP Addresses based on this container (like 172.18.0.20-172.18.0.50):`echo $'\n> '` " IPADDR_RANGE
echo "IP Address range is ${IPADDR_RANGE}"
echo

sed s/%ADDR_RANGE%/$IPADDR_RANGE/g \
   configs/metallb-config-template.yaml > configs/metallb-config-set.yaml

echo "Applying metallb-config-set.yaml with ip address range applied..."
kubectl apply -f configs/metallb-config-set.yaml

echo "Done! You should be ready to 'waypoint install -platform=kubernetes -accept-tos' on a local kubernetes!"
exit 0
