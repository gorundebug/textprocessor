#!/bin/bash

trap 'kill $(jobs -p)' SIGINT
./bin/wordsprocessor/wordsprocessor -config ./configs/wordsprocessor/config.yaml -values ./configs/wordsprocessor/values.yaml &
./bin/charsprocessor/charsprocessor -config ./configs/charsprocessor/config.yaml -values ./configs/charsprocessor/values.yaml &

wait