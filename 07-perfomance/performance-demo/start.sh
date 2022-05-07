#!/bin/bash

RUN_NAME=$1

kubectl create ns performance

helm install performance-test . --namespace= performance -f values.yaml --set run_name=$RUN_NAME