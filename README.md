# Stately

An unpretentious manager of files.


# Examples

``` shell
stately copy ./test -o tmp/
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


# Updating repos

``` shell
bazel run //:gazelle -- update-repos  -from_file=go.mod
```
