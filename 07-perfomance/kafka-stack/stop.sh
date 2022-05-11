#!/bin/bash

helm uninstall confluent-kafka --namespace confluent

helm uninstall confluent-operator --namespace confluent

kubectl delete -f data