#!/bin/bash

for INDEX in 0
do
kubectl delete -f data/data-0$INDEX.yaml
done