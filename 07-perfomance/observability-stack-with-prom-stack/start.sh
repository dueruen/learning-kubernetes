#!/bin/bash

RUN_NAME=$1

. create_data.sh first

helm dependency build
helm install observ . --create-namespace --namespace observability --set run_name=$RUN_NAME