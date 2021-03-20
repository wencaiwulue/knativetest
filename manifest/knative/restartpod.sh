#! /bin/bash

# shellcheck disable=SC2046
kubectl delete pods $(kubectl get pods -n test | grep webhook | awk -F' ' '{print$1}') -n test
