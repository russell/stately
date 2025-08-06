#!/usr/bin/env bats
# -*- mode: sh -*-
source "tests/common.sh"

setup_file() {
    export OUTPUT_DIR=$(mktemp -d)
    stately manifest -o "$OUTPUT_DIR" -i - < "tests/jsonnet-render-test.json"
}

teardown_file() {
    rm -rf "$OUTPUT_DIR"
}

@test "Test file examples/a.jsonnet" {
    test "$(file_type examples/a.jsonnet)" = "regular file"
    local FILE_CONTENTS=$(cat<<"EOF"
// a.jsonnet

local foo = {
  hello: 'world',
};

foo
EOF
)

    test "$(file_contents examples/a.jsonnet)" = "$FILE_CONTENTS"
}
