#!/bin/bash

CMD="cmd_build_bazel.sh"

bash ./version_inc.sh "./cmd/wordsprocessor/main.go"
bash ./version_inc.sh "./cmd/charsprocessor/main.go"


if ! bash build.sh ${CMD} "$1"; then
    echo "bash build.sh ${CMD} $1 failed"
    exit 1
fi

exit 0