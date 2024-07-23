# FB2EPUB converter

Command line converter from FB2 to EPUB. Supports transliteration of Cyrillic for header data.

```sh
$ go get github.com/reinventer/fb2epub
$ fb2epub -f source.fb2 -t destination.epub
```

## Limitations
* Can not convert from zipped FB2

## TODO
* Tests
* Support for zipped FB2
* To make correct table of contents