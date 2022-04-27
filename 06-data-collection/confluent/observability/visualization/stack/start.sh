#!/bin/bash

helm upgrade -f ./values.yaml --install observ grafana/loki-stack --namespace observability

kubectl get secret --namespace observability loki-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo