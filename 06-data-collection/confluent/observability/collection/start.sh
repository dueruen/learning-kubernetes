#!/bin/bash

echo "Create namespace"
kubectl create ns observability


echo "Start otel"
kubectl apply -f ./otel -n observability


echo "Start fluent-bit"
helm repo add fluent https://fluent.github.io/helm-charts
helm install -f ./logging/values.yaml fluent-bit fluent/fluent-bit --namespace observability 


echo "Start kube-state metrics"
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install kube-state-metrics prometheus-community/kube-state-metrics --version 4.7.0 --namespace observability
