#!/usr/bin/env bash

while : ; do
  kubectl logs --tail=1000 --follow=true $1 phylaxd
  sleep 1
done
