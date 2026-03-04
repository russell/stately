#!/usr/bin/env bats
# -*- mode: sh -*-
source "tests/common.sh"

@test "MergeSection creates file with managed block when file does not exist" {
    export OUTPUT_DIR=$(mktemp -d)
    stately manifest -o "$OUTPUT_DIR" -i - < "tests/merge-section-test.json"

    test "$(file_type managed-file)" = "regular file"
    run file_contents managed-file
    [ "$output" = "# BEGIN MANAGED
managed-line-1
managed-line-2
# END MANAGED" ]

    rm -rf "$OUTPUT_DIR"
}

@test "MergeSection appends managed block to existing file without markers" {
    export OUTPUT_DIR=$(mktemp -d)
    echo "user-content" > "$OUTPUT_DIR/managed-file"

    stately manifest -o "$OUTPUT_DIR" -i - < "tests/merge-section-test.json"

    run file_contents managed-file
    [ "$output" = "user-content

# BEGIN MANAGED
managed-line-1
managed-line-2
# END MANAGED" ]

    rm -rf "$OUTPUT_DIR"
}

@test "MergeSection replaces managed block in existing file with markers" {
    export OUTPUT_DIR=$(mktemp -d)
    cat > "$OUTPUT_DIR/managed-file" << 'EOF'
user-before
# BEGIN MANAGED
old-managed-content
# END MANAGED
user-after
EOF

    stately manifest -o "$OUTPUT_DIR" -i - < "tests/merge-section-test.json"

    run file_contents managed-file
    [ "$output" = "user-before
# BEGIN MANAGED
managed-line-1
managed-line-2
# END MANAGED
user-after" ]

    rm -rf "$OUTPUT_DIR"
}

@test "MergeSection cleanup removes only managed section, preserves user content" {
    export OUTPUT_DIR=$(mktemp -d)

    # First run: create the managed file
    stately manifest -s "$OUTPUT_DIR/.manifest.yml" -o "$OUTPUT_DIR" -i - < "tests/merge-section-test.json"

    # Add user content around the managed section
    cat > "$OUTPUT_DIR/managed-file" << 'EOF'
user-before
# BEGIN MANAGED
managed-line-1
managed-line-2
# END MANAGED
user-after
EOF

    # Second run: empty manifest triggers cleanup of the managed file
    stately manifest -s "$OUTPUT_DIR/.manifest.yml" -o "$OUTPUT_DIR" -i - < "tests/merge-section-empty-test.json"

    # File should still exist with only user content
    test -f "$OUTPUT_DIR/managed-file"
    run file_contents managed-file
    [ "$output" = "user-before
user-after" ]

    rm -rf "$OUTPUT_DIR"
}

@test "MergeSection cleanup deletes file if it contains only managed content" {
    export OUTPUT_DIR=$(mktemp -d)

    # First run: create the managed file
    stately manifest -s "$OUTPUT_DIR/.manifest.yml" -o "$OUTPUT_DIR" -i - < "tests/merge-section-test.json"

    # Second run: empty manifest triggers cleanup
    stately manifest -s "$OUTPUT_DIR/.manifest.yml" -o "$OUTPUT_DIR" -i - < "tests/merge-section-empty-test.json"

    # File should be deleted since it only had managed content
    test ! -f "$OUTPUT_DIR/managed-file"

    rm -rf "$OUTPUT_DIR"
}
