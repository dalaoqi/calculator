#!/bin/bash
DIR=$(dirname "$0")/..
pushd $DIR
CONTAINER="calculator"
DOCKERFILE="build/Dockerfile"

DOCKER_IMAGE=$CONTAINER

# Build docker
cmd="docker build -t $DOCKER_IMAGE -f $DOCKERFILE ."
echo "COMMAND:"$cmd
eval $cmd 

