#!/bin/bash

go mod tidy
trap 'kill $(jobs -p)' SIGINT
go run ./cmd/wordsprocessor/ -config ./configs/wordsprocessor/config.yaml -values ./configs/wordsprocessor/values.yaml &
go run ./cmd/charsprocessor/ -config ./configs/charsprocessor/config.yaml -values ./configs/charsprocessor/values.yaml &

wait