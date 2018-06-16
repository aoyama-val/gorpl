A command-line tool to find and replace strings in multiple files, written in Go.
It is designed for simplicity and friendly-output, and not a complete replacement for Joe Laffey's/Debian's rpl.
It resolves some drawbacks that perl/sed have (See below).

# Usage

```
Usage: rpl [options...] <from> <to> files...
  -i    ignore case
  -r    regular expression search
  -w    match whole word
```

```
$ cat normal
hoge
moge
hoge hoge
$ rpl hoge fuga normal
Search: /hoge/

Replace    normal    (3 matches)

1 files (replaced: 1 / no change: 0 / ignored: 0) Total 3 matches
$ cat normal
fuga
moge
fuga fuga
```


# Drawbacks of using Perl/sed for in-place string substitution

```
$ perl -pi -e 's/hoge/foo/g' *.c
```

- Symlinks are turned into regular files.
- Timestamps are updated even if no changes are made, no good for build tools like `make`.

`rpl` manages these points.
