#!/bin/bash

IT=
if [ "$GORUNDEBUG_TEST" != "1" ]; then
    IT='-it'
fi
docker run $IT --rm -p 9201:9201 -p 9202:9202 -p 9092:9092 -p 8080:8080 -p 9091:9091 -v "$(pwd)":/textprocessor -w /textprocessor -e GORUNDEBUG_TEST=$GORUNDEBUG_TEST -e CMD="cmd_run.sh" --name textprocessor_run textprocessor