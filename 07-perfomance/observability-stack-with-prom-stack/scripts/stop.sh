#!/bin/bash

helm uninstall observ . --namespace observability

kubectl delete -f /data