#!/bin/bash

RUN_NAME=$1
MESSAGE_SIZES=$2
MESSAGE_FREQUENCY=$3

if [[ -z "${MESSAGE_SIZES}" ]]; then
  echo "Using default message size"
  helm upgrade --install performance-test . --namespace performance --create-namespace \
    -f values.yaml \
    --set run_name=$RUN_NAME
else
  helm upgrade --install performance-test . --namespace performance --create-namespace \
    -f values.yaml \
    --set run_name=$RUN_NAME \
    --set producer.messageSize=$MESSAGE_SIZES \
    --set producer.messageFrequency=$MESSAGE_FREQUENCY
fi