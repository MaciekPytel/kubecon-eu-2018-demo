#!/bin/bash

# basically copy-paste from https://github.com/kubernetes-incubator/apiserver-builder/blob/master/docs/concepts/auth.md
PURPOSE="hpademo"
SERVICE_NAME=api
ALT_NAMES='"api.custom-metrics","api.custom-metrics.svc"'

openssl req -x509 -sha256 -new -nodes -days 365 -newkey rsa:2048 -keyout ${PURPOSE}-ca.key -out ${PURPOSE}-ca.crt -subj "/CN=ca"
echo '{"signing":{"default":{"expiry":"43800h","usages":["signing","key encipherment","'${PURPOSE}'"]}}}' > "${PURPOSE}-ca-config.json"
echo '{"CN":"'${SERVICE_NAME}'","hosts":['${ALT_NAMES}'],"key":{"algo":"rsa","size":2048}}' | cfssl gencert -ca=${PURPOSE}-ca.crt -ca-key=${PURPOSE}-ca.key -config=${PURPOSE}-ca-config.json - | cfssljson -bare apiserver

# write secret with key / cert
kubectl -n custom-metrics create secret tls cm-adapter-serving-certs --cert=apiserver.pem --key=apiserver-key.pem


# cleanup
rm ${PURPOSE}-ca*
rm apiserver*
