#!/bin/bash

. ./run.config

cd ../../application/performance-data
for value in $RUN_NAMES
do
echo $value

INPUT_PATH=/mnt/$VOLUME_FOLDER/data-00 FILE_NAME=data-$value FILE_EXTENSION=.json OUTPUT_PATH=/mnt/$VOLUME_FOLDER/data-00/output ./convert-to-csv

done