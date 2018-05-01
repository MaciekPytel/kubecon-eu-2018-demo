#!/bin/bash

kubectl delete -f hpa.yaml
kubectl delete -f load.yaml
kubectl delete -f ../adapter-config/custom-metrics-apiserver-deployment.yaml
kubectl delete -f ../cpu-consumer/deploy/server.yaml
kubectl create -f ../cpu-consumer/deploy/server.yaml
kubectl create -f load.yaml

STEP=0
