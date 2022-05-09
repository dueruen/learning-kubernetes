#!/bin/bash

. ./run.config

rm -rf /mnt/$VOLUME_FOLDER
mkdir /mnt/$VOLUME_FOLDER

if $WITH_CILIUM_METRICS
then

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

fi

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

##############################

if [ $count -eq 0 ]
then
   cd ../../observability-stack-with-prom-stack
   ./start.sh $VOLUME_FOLDER $HOST_NAME $value false

   cd ../kafka-stack
   ./start.sh $VOLUME_FOLDER $HOST_NAME

   echo "Wait for kafka to start Sleep 4m"
   for i in {1..4}
   do
      sleep 1m
      echo "$i m ..."
   done
   echo "Done waiting"
else
   cd ../observability-stack-with-prom-stack
   ./start.sh $VOLUME_FOLDER $HOST_NAME $value true
fi

cd ../performance-demo
./start.sh $value $ms $mf

echo "Round started - it will take $ROUND_RUN_TIME m"
for i in {1..$ROUND_RUN_TIME}
do
   sleep 1m
   echo "$i m ..."
done
echo "Round done"

##############################

count=$(( $count + 1 ))
MESSAGE_SIZE_OFFSET=$(( $MESSAGE_SIZE_OFFSET + 1 ))
if [ $MESSAGE_SIZE_OFFSET -eq $MESSAGE_SIZE_RUNS ]
then
MESSAGE_SIZE_OFFSET=0

fi

MESSAGE_FREQUENCY_OFFSET=$(( $count / $MESSAGE_SIZE_RUNS )) 
echo ""
done

echo "All done - but resources is not deleted"