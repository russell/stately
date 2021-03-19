load("@com_github_aignas_rules_shellcheck//:def.bzl", "shellcheck", "shellcheck_test")

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

sh_test(
    name = "stately_integration_copy",
    size = "small",
    srcs = ["copy_test.sh"],
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
