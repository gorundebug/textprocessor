#!/bin/bash

if ! make gen; then
    echo "make gen failed"
    exit 1
fi

if ! go mod tidy; then
    echo "go mod tidy"
    exit 1
fi

if ! bazel mod tidy; then
    echo "bazel mod tidy failed"
    exit 1
fi


if ! bazel build //...; then
    echo "bazel build //... failed"
    exit 1
fi

if ! bazel build //:install; then
    echo "bazel build //:install failed"
    exit 1
fi

exit 0