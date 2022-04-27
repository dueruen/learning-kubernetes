#!/bin/bash
echo "VOLUME_FOLDER: $1";
echo "HOST_NAME: $2";

if [[ $# -eq 0 ]] ; then
    echo 'Missing script argument'
    exit 0
fi

kubectl create namespace confluent

helm repo add confluentinc https://packages.confluent.io/helm
helm repo update

helm upgrade --install confluent-operator confluentinc/confluent-for-kubernetes --namespace confluent

VOLUME_FOLDER=$1 # volume_ams3_03
HOST_NAME=$2 # ubuntu-m-2vcpu-16gb-ams3-01

rm -rf data
mkdir data

for INDEX in 1 2 3
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

kubectl apply -f data -n confluent

# kubectl apply -f kafka -n confluent

(
  set -e
  kubectl wait deployment --namespace="cert-manager" --for="condition=Available" cert-manager-webhook cert-manager-cainjector cert-manager --timeout=3m
  kubectl wait pods --namespace="cert-manager" --for="condition=Ready" --all --timeout=3m
  kubectl wait apiservice --for="condition=Available" v1.cert-manager.io v1.acme.cert-manager.io --timeout=3m
  kubectl wait pods --namespace="confluent" --for="condition=Ready" --all --timeout=3m
  until kubectl get secret --namespace="cert-manager" cert-manager-webhook-ca 2> /dev/null ; do sleep 0.5 ; done
)