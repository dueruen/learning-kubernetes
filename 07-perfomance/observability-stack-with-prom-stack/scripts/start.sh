#!/bin/bash

VOLUME_FOLDER=$1
HOST_NAME=$2
RUN_NAME=$3
UPGRADE=$4

if $UPGRADE
then
  ./create_data_new_run.sh $VOLUME_FOLDER $RUN_NAME
else
  ./create_data.sh $VOLUME_FOLDER $HOST_NAME $RUN_NAME
fi

helm dependency build

helm upgrade --install observ . \
  --create-namespace --namespace observability \
  -f values.yaml \
  --set run_name=$RUN_NAME