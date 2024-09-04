bazel mod tidy

bazel build //...

bazel clean --expunge


make gen

make build

make run


docker build -t textprocessor .