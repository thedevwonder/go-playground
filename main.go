package main

import (
	sslTerminator "main/sslterminator"
)

func main () {
	sslTerminator.Helper("certs/cert.pem", "certs/key.pem")
}