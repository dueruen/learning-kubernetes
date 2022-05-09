#!/bin/bash

. ./run.config

count=0
MESSAGE_SIZE_OFFSET=0
MESSAGE_FREQUENCY_OFFSET=0

for value in $RUN_NAMES
do
echo $value

MULTI_SIZE=$((10 ** $MESSAGE_SIZE_OFFSET))
MULTI_FREQUENCY=$((10 ** $MESSAGE_FREQUENCY_OFFSET))

ms=$(( $MULTI_SIZE * $MESSAGE_SIZES_BASE ))
echo "Message size: " $ms

mf=$(( $MULTI_FREQUENCY * $MESSAGE_FREQUENCY_BASE ))
echo "Message frequency: " $mf

# echo $num
count=$(( $count + 1 ))
MESSAGE_SIZE_OFFSET=$(( $MESSAGE_SIZE_OFFSET + 1 ))
if [ $MESSAGE_SIZE_OFFSET -eq $MESSAGE_SIZE_RUNS ]
then
MESSAGE_SIZE_OFFSET=0
fi

MESSAGE_FREQUENCY_OFFSET=$(( $count / $MESSAGE_SIZE_RUNS )) 

echo ""
done