#!/bin/sh

set -m

cmd="/usr/src/app/server & /usr/src/app/client1 & /usr/src/app/client2 & /usr/src/app/client3;"
echo "COMMAND: exec "$cmd
exec $cmd
fg %1