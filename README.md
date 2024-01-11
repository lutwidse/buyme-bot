```bash
$ git clone -b temp-main-code https://github.com/lutwidse/buyme-bot.git
$ cd buyme-bot
$ go mod download
$ go build -o buyme cmd/main.go
```

In this temporary main code, the minimum required arguments that need to be configured are Discord and Proxy.

```bash
$ mv .env.example .env
```

Once you set that up, run it.

```bash
$ chmod +x buyme
$ ./buyme --debug true
```
