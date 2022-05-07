#!/bin/bash

RUN_NAME=$1

./create_data.sh $RUN_NAME

helm dependency build
helm install observ . --create-namespace --namespace observability -f values.yaml --set run_name=$RUN_NAME