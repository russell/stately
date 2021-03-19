#!/bin/bash

STATELY="$1"

OUTPUT_DIR=$(mktemp -d)

DIR="$( cd "$( dirname "$0" )" && pwd )"
$STATELY manifest -o "$OUTPUT_DIR" < "$DIR/test.json"

echo "===FILES==="
find "$OUTPUT_DIR"

file_type() {
    stat -c "%F" "$OUTPUT_DIR/$1"
}

set -xe
test "$(file_type c/foo3)" = "regular file"
test "$(file_type c/foo2)" = "regular file"
test "$(file_type c/foo1)" = "regular file"
test "$(file_type c/foo)" = "regular file"
test "$(file_type empty-file)" = "regular empty file"

rm -rf "$OUTPUT_DIR"
