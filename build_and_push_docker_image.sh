#!/bin/bash
set -e

IMAGE_TAG=$1
TIME=`date +%s`

# Create an Image from SCRATCH and push
docker build -t tchaudhry/mail-proxy-service:$IMAGE_TAG -f Dockerfile .
docker push tchaudhry/mail-proxy-service:$IMAGE_TAG
