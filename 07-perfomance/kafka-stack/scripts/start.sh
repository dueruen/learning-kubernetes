#!/bin/bash
VOLUME_FOLDER=$1
HOST_NAME=$2

kubectl create namespace confluent

./create_data.sh VOLUME_FOLDER HOST_NAME

helm repo add confluentinc https://packages.confluent.io/helm
helm repo update

helm upgrade --install confluent-operator confluentinc/confluent-for-kubernetes --namespace confluent

echo "Sleep 15s"
sleep 15s

helm install confluent-kafka . --namespace confluent -f values.yaml