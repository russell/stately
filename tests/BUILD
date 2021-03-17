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
