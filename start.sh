#!/bin/bash

export SCRIPTDIR="$(cd "$(dirname "$0")"; pwd -P)"

docker rm -v -f block_landing

docker run -d \
  --name block_landing \
  --restart=always \
  -p 9000:9000 \
  -v ${SCRIPTDIR}:/var/www \
  block_landing
