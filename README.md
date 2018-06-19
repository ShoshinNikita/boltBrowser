# boltBrowser

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/ShoshinNikita/boltBrowser/graphs/commit-activity)
[![Go Report Card](https://goreportcard.com/badge/github.com/ShoshinNikita/boltBrowser)](https://goreportcard.com/report/github.com/ShoshinNikita/boltBrowser)
[![Coverage Status](https://coveralls.io/repos/github/ShoshinNikita/boltBrowser/badge.svg?branch=add-travis)](https://coveralls.io/github/ShoshinNikita/boltBrowser?branch=master)
[![GitHub license](https://img.shields.io/github/license/ShoshinNikita/boltBrowser.svg)](https://github.com/ShoshinNikita/boltBrowser)
[![Github all releases](https://img.shields.io/github/downloads/ShoshinNikita/boltBrowser/total.svg)](https://github.com/ShoshinNikita/boltBrowser/releases/)

boltBrowser is a web-based explorer for BoltDB.

## Features

+ You can work with several databases in one time
+ You can visit nested buckets
+ Minimalistic and simple interface
+ Opportunity to search records and buckets using regex
+ Opportunity to edit a db (create, delete, modify buckets and records)

[Examples of using the program](Examples.md)

![1](stuff/screenshot.png)

## How to start

1. Run the program (you can download the latest release [here](https://github.com/ShoshinNikita/boltBrowser/releases))
1. Go to [localhost:500](http://localhost:500)
1. Open the list of databases
1. Add a database by pressing sign '+'
1. Enjoy!

## Settings

You can change mode of converting `[]byte`. Just change functions `ConvertKey(b []byte) string` (or `ConvertValue()`) in [src/converters/converter.go](src/converters/converters.go)

__Note__: function will be used for converting all keys (or values). So, if your keys (or values) were converted from either `string` or `uint` program will crash.

### Flags

Flag | Default | Description
---- | ------ | -------
`-port` | `:500` | port of the website
`-offset` | `100` | number of records on single page
`-debug` | `false` | mode of debugging
`-checkVer` | `true` | should program check a new version
`-writeMode` | `true` | can program edit databases
`-openBrowser` | `true` | should the program open a browser automatically
`-neatWindow` | `true` | should the program open a neat window

### Security

For preventing of js-injection program changes some symbols

Old symbol | New symbol
---------- | ----------
`<` | `❮`
`>` | `❯`
`"` | `＂`
`'` | `ߴ`

Scheme of work:

1. User sends a request
1. Program changes all new symbols to old (backend)
1. Program get info from a db
1. Program sends a response
1. Program changes all old symbols to new (frontend – function `SafeParse()`)

## Additional info

Initial work was undertaken on [Bitbucket](https://bitbucket.org/ShoshinNikita/boltbrowser).

## License

[MIT License](LICENSE)