#!/usr/bin/env bats
# -*- mode: sh -*-
source "tests/common.sh"

setup_file() {
    export OUTPUT_DIR=$(mktemp -d)
    export TEST_DIR=$(mktemp -d)

    newfile 'dir/file2' 'contents'
    newfile 'dir/file' 'contents1'
    newemptyfile 'empty-file'
    ln -s "$TEST_DIR/dir/file" "$TEST_DIR/ln"
    stately copy -L "--strip-prefix=$TEST_DIR" -o "$OUTPUT_DIR" "$TEST_DIR"
}

teardown_file() {
    rm -rf "$TEST_DIR"
    rm -rf "$OUTPUT_DIR"
}

@test "Staging files works" {
    stately -h
}

@test "Test file in subdir was copied" {
    test "$(file_type dir/file)" = "regular file"
}

@test "Test another file in a subdir was copied" {
    test "$(file_type dir/file2)" = "regular file"
}

@test "Test file was copied" {
    test "$(file_type empty-file)" = "regular empty file"
}

@test "Check symlinked file was copied" {
    test "$(file_type ln)" = "regular file"
}
