#!/bin/bash

go mod tidy
go build -gcflags="all=-N -l" -o ./bin/wordsprocessor/service1 ./cmd/wordsprocessor/main.go
go build -gcflags="all=-N -l" -o ./bin/charsprocessor/service1 ./cmd/charsprocessor/main.go


/go/bin/dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec "/textprocessor/bin/charsprocessor/charsprocessor" \
-- -config "/textprocessor/configs/charsprocessor/config.yaml" -values "/textprocessor/configs/charsprocessor/values.yaml"
/go/bin/dlv --listen=:40001 --headless=true --api-version=2 --accept-multiclient exec "/textprocessor/bin/wordsprocessor/wordsprocessor" \
-- -config "/textprocessor/configs/wordsprocessor/config.yaml" -values "/textprocessor/configs/wordsprocessor/values.yaml"