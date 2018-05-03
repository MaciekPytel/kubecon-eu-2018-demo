#!/bin/bash

STEP=${STEP:-0}

DEMO[0]="cat hpa.yaml"
DEMO[1]="kubectl create -f hpa.yaml"
DEMO[2]="cat patch1.yaml"
DEMO[3]='kubectl patch hpa.v2beta1.autoscaling sample-server --patch "$(cat patch1.yaml)"'
DEMO[4]="kubectl create -f ../adapter-config/custom-metrics-apiserver-deployment.yaml"
DEMO[5]="cat patch2.yaml"
DEMO[6]='kubectl patch hpa.v2beta1.autoscaling sample-server --patch "$(cat patch2.yaml)"'

BREAK[0]=false
BREAK[1]=true
BREAK[2]=false
BREAK[3]=true
BREAK[4]=true
BREAK[5]=false
BREAK[6]=true

for i in $(seq ${STEP} ${#DEMO[@]})
do
  echo
  echo "\$ ${DEMO[${i}]}"
  eval ${DEMO[${i}]}
  STEP=$((${i} + 1))
  if [ "${BREAK[${i}]}" = true ]
  then
    break
  fi
  sleep 1
done
