# boltBrowser

boltBrowser is a web-browser for BoltDB.

## Features

+ You can work with several databases in one time
+ You can visit nested buckets
+ Minimalistic and simple interface

## How to start

1. Run the program (you can download compiled [program](program/boltBrowser_v1.0.7z))
1. Go to [localhost:500](http://localhost:500)
1. Add db by pressing sign '+'
1. Enjoy!

## Settings

+ You can change mode of converting slice of byte. Just change functions ConvertKey() or ConvertValue() in [src/converters/converter.go](src/converters/converters.go)
	__Note__: function will be used for converting all keys (or values). So, if your keys (or values) were converted from both `string` and `uint` program will crash.
+ You can change port by flag `-port`

## Additional info

Initial work was undertaken on [Bitbucket](https://bitbucket.org/ShoshinNikita/boltbrowser).

## License

[MIT License](LICENSE)