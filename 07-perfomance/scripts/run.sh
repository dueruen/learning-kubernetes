#!/bin/bash

rm -rf /mnt/$VOLUME_FOLDER
mkdir /mnt/$VOLUME_FOLDER

helm upgrade cilium cilium/cilium --version 1.11.2 \
   --namespace kube-system \
   --reuse-values \
   --set prometheus.enabled=true \
   --set operator.prometheus.enabled=true \
   --set hubble.metrics.enabled="{dns,drop,tcp,flow,port-distribution,icmp,http}"
#   --set hubble.relay.enabled=true \
#   --set hubble.ui.enabled=true \
#   --set hubble.metrics.enabled="{dns,drop,tcp,flow,port-distribution,icmp,http}"

echo "Wait for Cilium to upgrade Sleep 30sec"
sleep 30s

cd ./observability-stack-with-prom-stack
./start.sh $VOLUME_FOLDER $HOST_NAME $RUN_NAME

cd ../kafka-stack
./start.sh $VOLUME_FOLDER $HOST_NAME

echo "Wait for kafka to start Sleep 4m"
sleep 1m
echo "1m ..."
sleep 1m
echo "2m ..."
sleep 1m
echo "3m ..."
sleep 1m
echo "4m - Done"

cd ../performance-demo
./start.sh $RUN_NAME

echo "Grafana password"
kubectl get secret --namespace observability observ-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo