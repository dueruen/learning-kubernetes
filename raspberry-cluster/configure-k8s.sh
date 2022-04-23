#!/bin/bash

echo "Create cluster"
sudo kubeadm init --pod-network-cidr=10.244.0.0/16

mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# Single node cluster
kubectl taint nodes --all node-role.kubernetes.io/master-

echo "Install helm"
wget https://get.helm.sh/helm-v3.8.2-linux-amd64.tar.gz
tar -zxvf helm-v3.8.2-linux-amd64.tar.gz
mv linux-amd64/helm /usr/local/bin/helm
rm helm-v3.8.2-linux-amd64.tar.gz

echo "Install DNI"
helm repo add cilium https://helm.cilium.io/

helm install cilium cilium/cilium --version 1.11.4 \
  --namespace kube-system