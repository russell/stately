load("@com_github_aignas_rules_shellcheck//:def.bzl", "shellcheck", "shellcheck_test")
load("@bazel_bats//:rules.bzl", "bats_test")

sh_library(
    name = "test_lib",
    data = [
        ":common.sh",
    ],
)

filegroup(
    name = "fixtures",
    srcs = [
        "test.json",
    ],
)

bats_test(
    name = "copy_test",
    srcs = ["copy_test.bats"],
    deps = ["//:stately", ":test_lib"],
)

bats_test(
    name = "manifest_test",
    srcs = ["manifest_test.bats"],
    deps = ["//:stately", ":test_lib", "fixtures"],
)

shellcheck_test(
    name = "shellcheck_test",
    data = glob(["*.sh"]),
    tags = ["lint"],
)
