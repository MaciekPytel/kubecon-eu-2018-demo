#!/bin/bash

watch "kubectl get hpa.v2beta1.autoscaling sample-server && echo && echo 'Deployment:' &&  kubectl get deployment sample-server"
