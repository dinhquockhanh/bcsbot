# BCS BOT

BCS Bot is a command line tool, auto get point in popcat game. This bot just use for some event days, so the source code look not pretty.

## Installation

install the [Go](https://go.dev/doc/install).

## Build

- linux:
```bash
env GOOS=linux go build -o exec-filename
```
## Usage
```bash
./exec-filename -h
Usage of ./exec-filename:
  -diff int
        diff number (default 555)
  -host string
        the host popcat server (default "popcat.lnquy.com")
  -username string
        username that registered (default "dqkhanh")
        
./exec-filename -u user1 -u user2
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)