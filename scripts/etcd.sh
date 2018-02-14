#!/usr/bin/env bash
set -e
name="flow-etcd"
version="v3.3.0"
docker kill $name || :
docker rm $name || :
docker run \
  -d \
  -p 2379:2379 \
  -p 2380:2380 \
  --name $name \
  quay.io/coreos/etcd:$version \
  etcd -listen-client-urls=http://0.0.0.0:2379


docker ps
