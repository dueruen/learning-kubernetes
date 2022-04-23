#!/bin/bash
VOLUME_FOLDER=volume_ams3_02
HOST_NAME=ubuntu-m-2vcpu-16gb-ams3-01

for INDEX in 1 2 3
do
mkdir /mnt/$VOLUME_FOLDER/data-0$INDEX
chmod 777 /mnt/$VOLUME_FOLDER/data-0$INDEX
done

helm dependency build

helm install cp-kafka . --namespace confluent