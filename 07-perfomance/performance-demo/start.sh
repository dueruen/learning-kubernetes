#!/bin/bash

RUN_NAME=$1

kubectl create ns performance

helm install performance-test . --set run_name=$RUN_NAME