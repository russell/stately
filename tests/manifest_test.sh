#!/bin/bash

STATELY=$1

OUTPUT_DIR=`mktemp -d`

DIR="$( cd "$( dirname "$0" )" && pwd )"
cat $DIR/test.json | $STATELY manifest -o $OUTPUT_DIR

echo "===FILES==="
find $OUTPUT_DIR

file_type() {
    stat -c "%F" $OUTPUT_DIR/$1
}

set -e
test "$(file_type c/foo3)" = "regular file"
test "$(file_type c/foo2)" = "regular file"
test "$(file_type c/foo1)" = "regular file"
test "$(file_type c/foo)" = "regular file"
test "$(file_type b)" = "regular empty file"

rm -rf $OUTPUT_DIR
