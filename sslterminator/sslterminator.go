package sslterminator

import (
	"fmt"
	"crypto/tls"
	"net"
	"io"
	"runtime"
	"time"
)

func Helper () {

	runtime.GOMAXPROCS(1)

	certPath := "certs/ecdsa_cert.pem"
	keyPath := "certs/ecdsa_key.pem"

	fmt.Println("Starting ssl terminator")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err!=nil {
		fmt.Println("error in loading cert")
		return
	}

	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	ls, err := tls.Listen("tcp", ":443", &tlsConfig)

	if err != nil {
		fmt.Println("error in creating tls server", err)
		return
	}
	fmt.Println("client is running on port: 443", ls.Addr())

	for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Println("error in accepting client connections")
			//resource released
			conn.Close()
			return
		}
		go sslTerminator(conn)
	}
}

func AcceptConnections(ls net.Listener) {
	for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Println("error in accepting client connections")
			return
		}

		t1 := time.Now()
		Handshake(conn)
		t2 := time.Now()
		fmt.Println(t2.Sub(t1), "time taken for handshake")
		conn.Close()
	}
}

func sslTerminator(clientConn net.Conn) {
	err := Handshake(clientConn)
	if err != nil {
		//resource released
		clientConn.Close()
		return
	}
	serverConn, err := createServer()

	if err != nil {
		fmt.Println("error in creating server connection", err)
		return
	}

	fmt.Println("server is running on port: 8080", serverConn.LocalAddr())

	//resource release already handled in tunnel
	go tunnel(clientConn, serverConn)
	go tunnel(serverConn, clientConn)
}

func Handshake(conn net.Conn) error {
	clientConn, ok := conn.(*tls.Conn)
	if ok {
		/* this handshake will probably throw error if invoked connection from the browser, 
		because the browser will not send the client hello message
		the go client will send the client hello message and then the handshake will be successful
		check the difference in remote addresses of the clients
		*/
		err := clientConn.Handshake()
		if err != nil {
			fmt.Println("error in handshake", err)
			return err
		}
		return nil
	}
	err := fmt.Errorf("error in handshake")
	return err
}

func createServer() (net.Conn, error) {
		serverConn, err := net.Dial("tcp", "127.0.0.1:8080")

		if err != nil {
			fmt.Println("error in connecting to server", err)
			return nil, err
		}

		fmt.Println("server is running on port: 8080", serverConn.LocalAddr())

		return serverConn, nil
}

func tunnel(from io.ReadWriteCloser, to io.ReadWriteCloser) {
	defer from.Close()
	defer to.Close()
	io.Copy(from, to)
	fmt.Println("copied data from", from, "to", to)
}