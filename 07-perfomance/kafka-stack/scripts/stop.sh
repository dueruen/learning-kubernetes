#!/bin/bash

helm uninstall confluent-kafka . --namespace confluent

helm uninstall confluent-operator confluentinc/confluent-for-kubernetes --namespace confluent