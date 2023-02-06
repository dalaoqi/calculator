#!/bin/bash
DIR=$(dirname "$0")/..
pushd $DIR
CONTAINER="calculator"
DOCKERFILE="build/Dockerfile"

DOCKER_IMAGE=$CONTAINER

# kill docker 
docker rm -f $CONTAINER

# detect the OS
platform=`uname`
if [[ "$platform" == 'Darwin' ]]; then
   timeZone=""
else
    timeZone="-v /etc/localtime:/etc/localtime"
fi

# run
cmd="docker run -d --restart=always --name $CONTAINER \
    -p 8081:8081 \
    $timeZone \
    $DOCKER_IMAGE"

echo "COMMAND:"$cmd
eval $cmd
RESULT=$?

popd

exit $RESULT
