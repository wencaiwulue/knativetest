#!/bin/bash

# shellcheck disable=SC2046
kubectl delete pods $(kubectl get pods -n test | grep test | awk -F' ' '{print$1}') -n test
