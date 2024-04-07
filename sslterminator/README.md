## SSL Terminator

Create an SSL certificate using OpenSSL

```
    openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
```
You'll get options to enter country code, city and some additional info.

Start an HTTP server on port: 8080

```
    openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes
```

To start ssl terminator, import `sslterminator` to the main.go file and call the Helper function

P.S - this is a learning project inspired from [https://github.com/cmpxchg16/go-sslterminator][https://github.com/cmpxchg16/go-sslterminator]
