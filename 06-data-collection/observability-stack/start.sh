#!/bin/bash

. create_data.sh

helm dependency build
helm install observ . --create-namespace --namespace observability