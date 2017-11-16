# FB2EPUB converter

Command line converter from FB2 to EPUB. Supports transliteration of Cyrillic for header data.

```sh
$ go get github.com/reinventer/fb2epub
$ fb2epub -f source.fb2 -t destination.epub
```

## Limitations
* Can not convert from zipped FB2
* The original fb2 should be in UTF-8 character encoding

## TODO
* Tests
* Support for zipped FB2
* Support for encodings other than UTF-8
* To make correct table of contents