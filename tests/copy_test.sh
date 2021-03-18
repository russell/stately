#!/bin/bash

STATELY=$1

OUTPUT_DIR=`mktemp -d`

DIR="$( cd "$( dirname "$0" )" && pwd )"

TEST_DIR=`mktemp -d`

newemptyfile() {
    local name=$1
    local contents=$2
    filepath=$TEST_DIR/$name
    mkdir -p $(dirname $filepath)
    touch $filepath
}

newfile() {
    local name=$1
    local contents=$2
    filepath=$TEST_DIR/$name
    mkdir -p $(dirname $filepath)
    echo '$contents' > $filepath
}

newfile 'c/foo3' 'foo'
newfile 'c/foo2' 'foo'
newfile 'c/foo1' 'foo'
newfile 'c/foo' 'foo'
newemptyfile 'b'
ln -s $TEST_DIR/c/foo3 $TEST_DIR/ln


echo "===INPUT FILES==="
find $TEST_DIR

echo "===RUNNING STATELY==="
$STATELY copy -L --strip-prefix=$TEST_DIR -o $OUTPUT_DIR $TEST_DIR

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

rm -rf $TEST_DIR
rm -rf $OUTPUT_DIR
