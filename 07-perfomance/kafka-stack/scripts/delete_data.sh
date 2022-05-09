#!/bin/bash

for INDEX in 1 2 3
do
kubectl delete -f data/data-0$INDEX.yaml
done