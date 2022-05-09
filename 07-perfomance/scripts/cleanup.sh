#!/bin/bash

echo "DELETE performance"
cd ../performance-demo
./stop.sh

echo "DELETE kafka"
cd ../kafka-stack
./stop.sh

echo "DELETE observability"
cd ../observability-stack-with-prom-stack
./stop.sh