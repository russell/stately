load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f2dcd210c7095febe54b804bb1cd3a58fe8435a909db2ec04e31542631cf715c",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "de69a09dc70417580aabf20a28619bb3ef60d038470c7cf8442fafcf627c21cb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

# Declare indirect dependencies
go_rules_dependencies()

go_register_toolchains(version = "1.18")

gazelle_dependencies()

#
# Shell Check
#
git_repository(
    name = "com_github_aignas_rules_shellcheck",
    commit = "94b231c8475f067c60f77459b2b54f4bcacc5e73",
    remote = "https://github.com/aignas/rules_shellcheck.git",
)

load("@com_github_aignas_rules_shellcheck//:deps.bzl", "shellcheck_dependencies")

shellcheck_dependencies()

#
# Bats
#
git_repository(
    name = "bazel_bats",
    remote = "https://github.com/filmil/bazel-bats",
    tag = "v0.29.1",
)

load("@bazel_bats//:deps.bzl", "bazel_bats_dependencies")

bazel_bats_dependencies(version = "v1.3.0")
