#!/usr/bin/env bats
# -*- mode: sh -*-
source "tests/common.sh"

setup_file() {
    export OUTPUT_DIR=$(mktemp -d)
    export TEST_DIR=$(mktemp -d)

    newfile 'dir/file2' 'contents2'
    newfile 'dir/file' 'contents'
    newfile 'file2' 'contents2'
    newfile 'file' 'contents'
    stately copy -L "--strip-prefix=$TEST_DIR" -o "$OUTPUT_DIR" "$TEST_DIR"
}

teardown_file() {
    rm -rf "$TEST_DIR"
    rm -rf "$OUTPUT_DIR"
}

@test "Staging files works" {
    stately -h
}

@test "Test removing a file from the sourc subdir will be reflected" {
    test "$(file_type dir/file2)" = "regular file"
    rm "$TEST_DIR/dir/file2"
    stately copy -L "--strip-prefix=$TEST_DIR" -o "$OUTPUT_DIR" "$TEST_DIR"
    test ! -f "$OUTPUT_DIR/dir/file2"
}

@test "Test another file in a subdir was copied" {
    test "$(file_type file2)" = "regular file"
    rm "$TEST_DIR/file2"
    stately copy -L "--strip-prefix=$TEST_DIR" -o "$OUTPUT_DIR" "$TEST_DIR"
    test ! -f "$OUTPUT_DIR/file2"
}
