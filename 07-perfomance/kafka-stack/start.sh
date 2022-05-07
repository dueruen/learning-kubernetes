#!/bin/bash
kubectl create namespace confluent

. create_data.sh

helm repo add confluentinc https://packages.confluent.io/helm
helm repo update

helm upgrade --install confluent-operator confluentinc/confluent-for-kubernetes --namespace confluent

echo "Sleep 15s"
sleep 15s

helm install confluent-kafka . --namespace confluent