#!/bin/sh
echo "Output file name: $1"
echo "Running..."
echo "Memory(MB)-Used Memory(MB)-Total Memory-Used-Percent(%) CPU-Load(%) Date" > $1

while :
do
  n=$(free -m | awk 'NR==2{printf "%s %s %.2f\n", $3,$2,$3*100/$2}')
  n="$n $(top -bn1 | grep load | awk '{printf "%.2f\n", $(NF-2)}')"
  n="$n $(date +"%Y/%m/%d:%H-%M-%S-%N")"
  echo $n >> $1
  sleep 0.25;
done