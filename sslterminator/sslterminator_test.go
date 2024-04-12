package sslterminator

import (
	"testing"
	"runtime"
	"crypto/tls"
	"os/exec"
	"net"
	"fmt"
)

var listener net.Listener

func benchmarkTLSConnection(t *testing.B, certPath string, keyPath string, concurrency int) {

	if listener == nil {
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err!=nil {
			t.Fatal("error in loading cert")
			return
		}
	
		tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		ls, err := tls.Listen("tcp", ":443", &tlsConfig)

		listener = ls
	
		if err != nil {
			t.Fatal("error in creating tls server", err)
			return
		}
	}

	listenerCloseCh := make(chan net.Listener)
	
	go AcceptConnections(listener)
	// fmt.Println("client is running on port: 443", ls.Addr())
	if concurrency == 1 {
		callHTTPS(listenerCloseCh)
	}

	<-listenerCloseCh
}

func callHTTPS(ch chan net.Listener) {
	arg0 := "openssl"
	arg1 := "s_client"
	arg2 := "localhost:443"
	cmd := exec.Command(arg0, arg1, arg2)
	_, sslErr := cmd.Output()
	if sslErr != nil {
		fmt.Println("error in ssl client", sslErr)
	}
	close(ch)
}

func BenchmarkUsingECDSAKeys(t *testing.B) {
	runtime.GOMAXPROCS(1)
	for i := 0; i < t.N; i++ {
		benchmarkTLSConnection(t, "certs/ecdsa_cert.pem", "certs/ecdsa_key.pem", 1)
	}
}

func BenchmarkUsingRSAKeys(t *testing.B) {
	runtime.GOMAXPROCS(1)
	for i := 0; i < t.N; i++ {
		benchmarkTLSConnection(t, "certs/rsa_cert.pem", "certs/rsa_key.pem", 1)
	}
}

