#!/usr/bin/env bats
# -*- mode: sh -*-

newemptyfile() {
    local name="$1"
    local contents="$2"
    filepath="$TEST_DIR/$name"
    mkdir -p "$(dirname "$filepath")"
    touch "$filepath"
}

newfile() {
    local name="$1"
    local contents="$2"
    filepath="$TEST_DIR/$name"
    mkdir -p "$(dirname "$filepath")"
    echo "$contents" > "$filepath"
}

file_type() {
    stat -c "%F" "$OUTPUT_DIR/$1"
}

set_file_mode() {
    local name="$1"
    local mode="$2"
    chmod "$mode" "$TEST_DIR/$name"
}

file_mode() {
    stat -c "%a" "$OUTPUT_DIR/$1"
}
