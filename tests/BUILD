load("@bazel_bats//:rules.bzl", "bats_test")

# TODO: Fix shellcheck integration with Bzlmod
# load("@rules_shellcheck//:def.bzl", "shellcheck", "shellcheck_test")

sh_library(
    name = "test_lib",
    data = [
        ":common.sh",
    ],
)

filegroup(
    name = "fixtures",
    srcs = [
        "dhall-render-test.json",
        "test.json",
    ],
)

bats_test(
    name = "copy_test",
    srcs = ["copy_test.bats"],
    deps = [
        ":test_lib",
        "//:stately",
    ],
)

bats_test(
    name = "copy_filemode_test",
    srcs = ["copy_filemode_test.bats"],
    deps = [
        ":test_lib",
        "//:stately",
    ],
)

bats_test(
    name = "copy_update_test",
    srcs = ["copy_update_test.bats"],
    deps = [
        ":test_lib",
        "//:stately",
    ],
)

bats_test(
    name = "manifest_test",
    srcs = ["manifest_test.bats"],
    data = [":fixtures"],
    deps = [
        ":test_lib",
        "//:stately",
    ],
)

bats_test(
    name = "manifest_dhall_test",
    srcs = ["manifest_dhall_test.bats"],
    data = [":fixtures"],
    deps = [
        ":test_lib",
        "//:stately",
    ],
)

# shellcheck_test(
#     name = "shellcheck_test",
#     data = glob(["*.sh"]),
#     tags = ["lint"],
# )
