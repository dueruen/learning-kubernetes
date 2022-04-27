#!/bin/bash
# https://www.jaegertracing.io/docs/1.33/operator/

kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.6.3/cert-manager.yaml

(
  set -e
  kubectl wait deployment --namespace="cert-manager" --for="condition=Available" cert-manager-webhook cert-manager-cainjector cert-manager --timeout=3m
  kubectl wait pods --namespace="cert-manager" --for="condition=Ready" --all --timeout=3m
  kubectl wait apiservice --for="condition=Available" v1.cert-manager.io v1.acme.cert-manager.io --timeout=3m
  until kubectl get secret --namespace="cert-manager" cert-manager-webhook-ca 2> /dev/null ; do sleep 0.5 ; done
)

sleep 10s

kubectl create namespace observability
kubectl create -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.33.0/jaeger-operator.yaml -n observability

cat > jaeger.yaml << EOF
apiVersion: jaegertracing.io/v1
kind: Jaeger
metadata:
  name: jaeger-default
  namespace: jaeger
spec:
  strategy: allInOne
  storage:
    type: memory
    options:
      memory:
        max-traces: 100000
  ingress:
    enabled: false
  annotations:
    scheduler.alpha.kubernetes.io/critical-pod: ""
EOF
kubectl apply -f jaeger.yaml -n observability