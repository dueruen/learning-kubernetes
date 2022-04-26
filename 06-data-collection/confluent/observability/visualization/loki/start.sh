#!/bin/bash
echo "VOLUME_FOLDER: $1";
echo "HOST_NAME: $2";

if [[ $# -eq 0 ]] ; then
    echo 'Missing script argument'
    exit 0
fi

VOLUME_FOLDER=$1 # volume_ams3_03
HOST_NAME=$2 # ubuntu-m-2vcpu-16gb-ams3-01

rm -rf data
mkdir data

for INDEX in 4 5
do
rm -rf /mnt/$VOLUME_FOLDER/data-0$INDEX
mkdir /mnt/$VOLUME_FOLDER/data-0$INDEX
chmod 777 /mnt/$VOLUME_FOLDER/data-0$INDEX
cat > data/data-0$INDEX.yaml << EOF
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: data-0$INDEX
spec:
  capacity:
    storage: 10Gi
  volumeMode: Filesystem
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Delete
  local:
     path: /mnt/$VOLUME_FOLDER/data-0$INDEX # Must be a path on the worker node
  nodeAffinity:
    required:
      nodeSelectorTerms:
      - matchExpressions:
        - key: kubernetes.io/hostname
          operator: In
          values:
          - $HOST_NAME
---
EOF
done

helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm upgrade --install vis --namespace=loki grafana/loki-simple-scalable
