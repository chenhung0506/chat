#!/bin/bash
REPO=harbor.linch.live/go
CONTAINER=chat
# TAG=$(git rev-parse --short HEAD)-$(date '+%Y%m%d-%H%M') 
TAG="latest"
DOCKER_IMAGE=$REPO/$CONTAINER:$TAG
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# BUILDROOT=$DIR/..
BUILDROOT=$DIR/
echo $BUILDROOT

# Windows
# SELF_IP=`ifconfig | grep -A 1 eth0 | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $2}' | head -n1`;
SELF_IP=`ip addr | grep -A 1 eth0 | grep "inet " | grep -Fv 127.0.0.1 | awk '{print $2}' | awk -F '/' '{print $1}'`;

echo $SELF_IP

API_IP=$SELF_IP:3002

cmd="docker build -t $DOCKER_IMAGE -f $DIR/DockerFile $BUILDROOT"

echo $cmd
eval $cmd


echo $DOCKER_IMAGE

docker rm -f chat
docker run --name chat --restart always -p 3002:3002 -d $DOCKER_IMAGE