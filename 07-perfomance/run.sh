#!/bin/bash

RUN_NAME="first"

helm upgrade cilium cilium/cilium --version 1.11.2 \
   --namespace kube-system \
   --reuse-values \
   --set prometheus.enabled=true \
   --set operator.prometheus.enabled=true \
  #  --set hubble.relay.enabled=true \
  #  --set hubble.ui.enabled=true \
   --set hubble.metrics.enabled="{dns,drop,tcp,flow,port-distribution,icmp,http}"

echo "Wait for Cilium to upgrade Sleep 30sec"
sleep 30s

./observability-stack-with-prom-stack/start.sh $RUN_NAME

./kafka-stack/start.sh

echo "Wait for kafka to start Sleep 4m"
sleep 4m

./performance-demo/start.sh $RUN_NAME