#!/bin/bash

increment_version() {
    local version="$1"
    IFS='.' read -r major minor patch <<< "${version#v}"
    patch=$((patch + 1))
    echo "v$major.$minor.$patch"
}

file_path="$1"

content=$(cat "$file_path")
if [[ $content =~ build_version\ =\ \"([^\"]*)\" ]]; then
    build_version="${BASH_REMATCH[1]}"
    new_version=$(increment_version "$build_version")
    updated_content=$(echo "$content" | sed "s/$build_version/$new_version/")
    echo "$updated_content" > "$file_path"
    echo "Updated version: $new_version"
else
    echo "Warning: build_version variable not found"
fi

if ! commit_hash=$(git rev-parse HEAD 2>/dev/null); then
    echo "Warning: commit not found."
else
    echo "${commit_hash}"
    content=$(cat "$file_path")
    if [[ $content =~ build_commit\ =\ \"([^\"]*)\" ]]; then
        updated_content=$(echo "$content" | sed "s/build_commit = \"[^\"]*\"/build_commit = \"$commit_hash\"/")
        echo "$updated_content" > "$file_path"
        echo "Updated commit: $commit_hash"
    else
        echo "Warning: build_commit variable not found"
    fi
fi