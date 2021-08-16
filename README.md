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

This project implement the same input standard that is used in
[dhall-render](https://github.com/timbertson/dhall-render), it's not
really a standard but I already had some tools using it.

# Modes

Stately supports 2 modes of operation, the first is `copy`.  This will
copy a source directory structure to a destination.  The second on is
the `manifest`, in this mode an input JSON file is converted into a
filestructure.  This mode is mostly useful if you are using a language
like Jsonnet, and you want to do more complicated templating like
outputting executables.

## Copy

Running the copy command will effectively copy your files from one
location to another.

``` shell
$ stately copy ./test -o tmp/
2021-08-16T21:34:35.216+0200	DEBUG	actions/copy.go:124	Copying: test/file2.txt -> tmp/test/file2.txt
2021-08-16T21:34:35.217+0200	DEBUG	actions/copy.go:124	Copying: test/foo/file1.txt -> tmp/test/foo/file1.txt2```

If you use find you can see that the same files exist at the destination.

``` shell
$ find tmp
tmp/
tmp/test
tmp/test/file2.txt
tmp/test/foo
tmp/test/foo/file1.txt
```

But Stately keeps a record of what files are copied. It stores these
in the `.stately-files.yaml` file. You can optionally specify an
alternative file with `--state-file` or use `--name` to specify an
alternative target name.  The default target is `default` as can be
seen in this file.

``` yaml
apiVersion: simopolis.xyz/v1alpha1
kind: StateConfig
target:
  default:
    files:
    - Path: tmp/test/file2.txt
    - Path: tmp/test/foo/file1.txt
```

If we remove one of the files and copy again, then Stately will remove
the file that no longer exists at the source.

``` shell
$ rm test/foo/file1.txt
$ bazel-bin/stately_/stately copy ./test -o tmp/
2021-08-16T21:36:37.029+0200	DEBUG	actions/copy.go:124	Copying: test/file2.txt -> tmp/test/file2.txt
2021-08-16T21:36:37.029+0200	DEBUG	config/state_file.go:136	Deleting: tmp/test/foo/file1.txt
$ find tmp/
tmp/
tmp/test
tmp/test/file2.txt
tmp/test/foo
```

## Manifest

Manifst works with a JSON file like this `test.json` one.

``` json
{
    "files": {
        "foo1/file3.json": {
            "executable": false,
            "contents": "{\"key\": \"value\"}",
            "format": "JSON",
            "install": "Write"
        }
    }
}
```

When it's passed into stately, it will manifest the file into the target directory.

``` shell
$ bazel-bin/stately_/stately manifest --input test.json --output-dir tmp
2021-08-16T21:48:41.750+0200	DEBUG	actions/manifest.go:82	Manifesting file: tmp/foo1/file3.json
```

Stately will also create a state file `.stately-files.yaml` much like it did for copy.

``` yaml
apiVersion: simopolis.xyz/v1alpha1
kind: StateConfig
target:
  default:
    files:
    - Path: tmp/foo1/file3.json
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
bazel test ...
```
