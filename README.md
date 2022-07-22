# BCS BOT

BCS Bot is a command line tool, auto get point in popcat game. This bot just use for some event days, so the source code look not pretty.

## Installation

sudo wget -qO - https://github.com/dinhquockhanh/bcsbot/releases/download/v0.0.2/bscbot-v0.0.2-linux-amd64.tar.gz | tar zxv && chmod +x bscbot

## Build
install the [Go](https://go.dev/doc/install).

- linux:
```bash
env GOOS=linux go build -o exec-filename
```
## Usage
```bash
./exec-filename -h
Usage of ./exec-filename:
  -diff int
        the diff number, if long time you don't receive the point, plz try decrease the diff number, ex: -diff=100 (default 398)
  -host string
        the host popcat server (default "popcat.lnquy.com")
  -max int
        the max connections, if max = 5, you don't have connection for browser... (default 4)
  -rps int
        the request per second (default 1)
  -token string
        the team's password
  -u value
        list users name, ex: -u=user1 -u=user2

        
./exec-filename -u user1 -u user2 -token SDHFKASDJHFIUASDFKJH
 
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)