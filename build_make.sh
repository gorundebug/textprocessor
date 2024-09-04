#!/bin/bash

CMD="cmd_build_make.sh"

if ! bash build.sh ${CMD} "$1"; then
    echo "bash build.sh ${CMD} $1 failed"
    exit 1
fi

exit 0