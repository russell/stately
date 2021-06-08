# Stately

An unpretentious manager of files. When using templates to generate
config files, it's often the case that you don't want to have to
remember to remove old files that have been removed from the templates
output.

This command provides a way to copy files from one directory to
another recording the files that were copied so that should the next
time you generate the files there are less, it will remove any files
that were removed.

There are 2 modes, one is a simple copy, where the content of folder
is copied to another folder, recording the files that were copied.

The other takes a JSON input and will write it out as files.  This is
to support more complex use cases where you want to do symlinks or
other permissions, or headers.

# Examples

``` shell
$ stately copy ./test -o tmp/
2021-03-14T15:41:06.334+0100	DEBUG	actions/copy.go:58	Copying: test/foo
2021-03-14T15:41:06.334+0100	DEBUG	actions/copy.go:58	Copying: test/foo/c.txt
2021-03-14T15:41:06.334+0100	DEBUG	actions/copy.go:58	Copying: test/t.text
```

``` shell
$ find tmp
tmp
tmp/foo
tmp/foo/c.txt
tmp/t.text
```


# Contributing

## Updating dependencies

``` shell
bazel run //:gazelle -- update
bazel run //:gazelle -- update-repos  -from_file=go.mod
```

## Building

``` shell
bazel build //:stately
```

## Tests

``` shell
bazel test //:...
```
