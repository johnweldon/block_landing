#!/bin/bash

NAME=block_landing
IMAGE=block_landing
HOST="block.example.com"
SYSLOGHOST=169.254.10.10
PORT=9000

docker rm -v -f ${NAME}
docker run -d \
    --name ${NAME} \
    --restart=always \
    -e BLOCK_PORT=8999 \
    -e BLOCK_SYSLOG_HOST=${SYSLOGHOST} \
    -h ${HOST} \
    -p ${PORT}:8999 \
  ${IMAGE}

