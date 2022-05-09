#!/bin/bash

VOLUME_FOLDER=$1
HOST_NAME=$2
RUN_NAME=$3

rm -rf ./data
mkdir ./data

for INDEX in 0
do
rm -rf /mnt/$VOLUME_FOLDER/data-0$INDEX
mkdir /mnt/$VOLUME_FOLDER/data-0$INDEX
chmod 777 /mnt/$VOLUME_FOLDER/data-0$INDEX
echo "" > /mnt/$VOLUME_FOLDER/data-0$INDEX/data-$RUN_NAME.json
chmod 777 /mnt/$VOLUME_FOLDER/data-0$INDEX/data-$RUN_NAME.json

cat > data/data-0$INDEX.yaml << EOF
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: data-0$INDEX
spec:
  capacity:
    storage: 30Gi
  volumeMode: Filesystem
  accessModes:
  - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  storageClassName: my-storage-class
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

kubectl apply -f data