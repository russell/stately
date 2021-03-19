load("@com_github_aignas_rules_shellcheck//:def.bzl", "shellcheck", "shellcheck_test")
load("@bazel_bats//:rules.bzl", "bats_test")

sh_library(
    name = "test_lib",
    data = [
        ":common.sh",
    ],
)

bats_test(
    name = "copy_test",
    srcs = ["copy_test.bats"],
    deps = ["//:stately", ":test_lib"],
)

sh_test(
    name = "stately_integration_manifest",
    size = "small",
    srcs = ["manifest_test.sh"],
    args = ["$(location //:stately)"],
    data = [
        "test.json",
        "//:stately",
    ],
)

shellcheck_test(
    name = "shellcheck_test",
    data = glob(["*.sh"]),
    tags = ["lint"],
)
