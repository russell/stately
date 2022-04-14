load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/russell/stately
gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)

go_library(
    name = "stately_lib",
    srcs = ["main.go"],
    importpath = "github.com/russell/stately",
    visibility = ["//visibility:private"],
    deps = ["//cmd"],
)

go_binary(
    name = "stately",
    embed = [":stately_lib"],
    visibility = ["//visibility:public"],
)
