#!/usr/bin/env bats
# -*- mode: sh -*-
source "tests/common.sh"

setup_file() {
    export OUTPUT_DIR=$(mktemp -d)
    stately manifest -o "$OUTPUT_DIR" < "tests/test.json"
}

teardown_file() {
    rm -rf "$TEST_DIR"
    rm -rf "$OUTPUT_DIR"
}

@test "Staging files works" {
    stately -h
}

@test "Test file foo3 in subdir was created" {
    test "$(file_type c/foo3)" = "regular file"
}

@test "Test file foo2 in a subdir was created" {
    test "$(file_type c/foo2)" = "regular file"
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
