#!/bin/bash

VOLUME_FOLDER=$1
RUN_NAME=$2

for INDEX in 0
do
echo "" > /mnt/$VOLUME_FOLDER/data-0$INDEX/data-$RUN_NAME.json
chmod 777 /mnt/$VOLUME_FOLDER/data-0$INDEX/data-$RUN_NAME.json

done