#!/usr/bin/env bats
# -*- mode: sh -*-
source "tests/common.sh"

setup_file() {
    export OUTPUT_DIR=$(mktemp -d)
    export TEST_DIR=$(mktemp -d)

    newfile 'dir/file2' 'contents'
    set_file_mode 'dir/file2' 400
    newfile 'dir/file' 'contents1'
    set_file_mode 'dir/file2' 550
}

teardown_file() {
    rm -rf "$TEST_DIR"
    rm -rf "$OUTPUT_DIR"
}

@test "Test file in subdir was copied" {
    echo Original Mode: $(stat -c "%a" "$TEST_DIR/dir/file")
    test $(stat -c "%a" "$TEST_DIR/dir/file") = "644"
    stately copy -L "--strip-prefix=$TEST_DIR" -o "$OUTPUT_DIR" "$TEST_DIR"
    test "$(file_type dir/file)" = "regular file"
    echo Final Mode: $(file_mode dir/file)
    test "$(file_mode dir/file)" = "644"
}

@test "Test another file in a subdir was copied" {
    echo Original Mode: $(stat -c "%a" "$TEST_DIR/dir/file2")
    test $(stat -c "%a" "$TEST_DIR/dir/file") = "644"
    stately copy -L "--strip-prefix=$TEST_DIR" -o "$OUTPUT_DIR" "$TEST_DIR"
    test "$(file_type dir/file2)" = "regular file"
    echo Final Mode: $(file_mode dir/file2)
    test "$(file_mode dir/file2)" = "750"
}


@test "Test file in subdir was copied and permissions overridden" {
    echo Original Mode: $(stat -c "%a" "$TEST_DIR/dir/file")
    test $(stat -c "%a" "$TEST_DIR/dir/file") = "644"
    stately copy -L "--strip-prefix=$TEST_DIR" -o "$OUTPUT_DIR" --file-mode 777 "$TEST_DIR"
    test "$(file_type dir/file)" = "regular file"
    echo Final Mode: $(file_mode dir/file)
    test "$(file_mode dir/file)" = "777"
}

@test "Test another file in a subdir was copied and permissions overridden" {
    echo Original Mode: $(stat -c "%a" "$TEST_DIR/dir/file2")
    test $(stat -c "%a" "$TEST_DIR/dir/file") = "644"
    stately copy -L "--strip-prefix=$TEST_DIR" -o "$OUTPUT_DIR" --file-mode 0777 "$TEST_DIR"
    test "$(file_type dir/file2)" = "regular file"
    echo Final Mode: $(file_mode dir/file2)
    test "$(file_mode dir/file2)" = "777"
}
