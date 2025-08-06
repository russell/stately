#!/usr/bin/env bats
# -*- mode: sh -*-
source "tests/common.sh"

setup_file() {
    export OUTPUT_DIR=$(mktemp -d)
    stately manifest -o "$OUTPUT_DIR" -i - < "tests/test.json"
}

teardown_file() {
    rm -rf "$OUTPUT_DIR"
}

@test "Staging files works" {
    stately -h
}

@test "Test jsonnet-file was written" {
    test "$(file_type jsonnet/file)" = "regular file"
}

@test "Test file c/foo3 in subdir was created" {
    test "$(file_type c/foo3)" = "regular file"
}

@test "Test file c/foo2 in a subdir was created" {
    test "$(file_type c/foo2)" = "regular file"
}

@test "Test file c/foo2 is executable" {
    mode=$(stat -c "%a" "$OUTPUT_DIR/c/foo2")
    [ "$mode" = "755" ]
}

@test "Test file c/foo3 is not executable" {
    mode=$(stat -c "%a" "$OUTPUT_DIR/c/foo3")
    [ "$mode" = "644" ]
}

@test "Test file foo1 in a subdir was created" {
    test "$(file_type c/foo1)" = "regular file"
}

@test "Test file foo in a subdir was created" {
    test "$(file_type c/foo)" = "regular file"
}

@test "Empty file was copied" {
    test "$(file_type empty-file)" = "regular empty file"
}
