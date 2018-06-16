A command-line tool to find and replace strings in multiple files, written in Go.
It is designed for simplicity and friendly-output, and not a complete replacement for Joe Laffey's/Debian's rpl.

# Usage

```
Usage: ./rpl [OPTIONS] <from> <to> files...
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
