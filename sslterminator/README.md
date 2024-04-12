## SSL Terminator

Create an SSL certificate using RSA 2048

```
    openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
```

using ecdsa certificate which is lighter in computation

```
    go run `go env GOROOT`/src/crypto/tls/generate_cert.go --host=localhost --ecdsa-curve=P256
```
You'll get options to enter country code, city and some additional info.

Start an HTTP server on port: 8080

```
    while true; do { echo -e 'HTTP/1.1 200 OK\r\n'; } | nc -l 8080; done
```

To start ssl terminator, import `sslterminator` to the main.go file and call the Helper function

P.S - this is a learning project inspired by [https://github.com/cmpxchg16/go-sslterminator](https://github.com/cmpxchg16/go-sslterminator)
