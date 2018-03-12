# jig

Minimal Go/TOML based task runner.

## Install

You can download a binary from [release](https://github.com/tanksuzuki/jig/releases).
Distributions are available for the...

* Mac OS X
* Linux

## Usage

```
$ jig -h

Usage:
  jig [<option>] <script_name> [<script_arg>...]

Application Options:
  -c, --config=path    Config file to load (default: ~/jig.toml) [$JIG_CONFIG]
  -h, --help           Help for jig or Print your script usage
      --version        Print version information and quit
```

### Listing scripts

```
$ jig

edit             Edit ~/jig.toml with vim
example <arg>    This is an example script
sync             Synchronize jig.toml with GitHub Gist
```

### Script help

```
$ jig -h example

This is an example script

Usage:
  jig example <arg>
```

### Running scripts

For example, suspose you have the following script:

```
[[script]]
name = "example"
args = "<arg>"
description = "This is an example script"
exec = '''
echo '$#': $#
echo '$0': $0
echo '$1': $1
'''
```

You can use following command to run `example` script.
If you specify arguments after script name, it will be passed to script.

```
$ jig example arg1

$#: 1
$0: example
$1: arg1
```

The exit code of jig is the same as the exit code of the script.

## Configuration

By default, `~/jig.toml` is loaded.
If you want change the path, you can use `-c` option or set `$JIG_CONFIG`.

### Example

```
[[script]]
name = "example"
args = "<arg>" # Optional
description = "This is an example script" # Optional
exec = '''
echo '$#': $#
echo '$0': $0
echo '$1': $1
'''

[[script]]
name = "edit"
description = """
Edit ~/jig.toml with vim
If this field is multiline,
the second and subsequent lines will be displayed in the script help.
"""
exec = "vim ~/.jig.toml"
```

```
$ jig -h edit

Edit ~/jig.toml with vim

Usage:
  jig edit

If this field is multiline,
the second and subsequent lines will be displayed in the script help.
```

## License

MIT

## Author

[Asuka Suzuki](https://github.com/tanksuzuki)
