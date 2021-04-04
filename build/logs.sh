#!/bin/bash

IMAGE=$1:latest

# restart the pods
for POD in $(kubectl get pods -n test | awk -F ' ' '{print$1}' | grep -v NAME); do
  # shellcheck disable=SC2046
  if [ "$IMAGE" == $(kubectl get pods "$POD" -n test -o jsonpath='{.spec.containers[0].image}') ]; then
    kubectl wait --for=condition=Ready pod/"$POD" -n test
    kubectl logs pods/"$POD" -n test -f
    break
  fi
done
