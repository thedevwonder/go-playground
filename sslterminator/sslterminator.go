package sslterminator

import (
	"fmt"
	"crypto/tls"
	"net"
	"io"
)

func Helper (certPath, keyPath string) {

	fmt.Println("Starting ssl terminator")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err!=nil {
		fmt.Println("error in loading cert")
		return
	}

	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	client, err := tls.Listen("tcp", ":443", &tlsConfig)

	if err != nil {
		fmt.Println("error in creating tls server", err)
		return
	}
	fmt.Println("client is running on port: 443", client.Addr())

	for {
		conn, err := client.Accept()
		if err != nil {
			fmt.Println("error in accepting client connections")
			return
		}
		fmt.Println("client is running on port: 443", conn.LocalAddr(), conn.RemoteAddr())
		
		go sslTerminator(conn)
	}
}

func sslTerminator(clientConn net.Conn) {
	err := decrypt(clientConn)
	if err != nil {
		return
	}
	serverConn, err := createServer()

	if err != nil {
		fmt.Println("error in creating server connection", err)
		return
	}

	fmt.Println("server is running on port: 8080", serverConn.LocalAddr())

	go tunnel(clientConn, serverConn)
	go tunnel(serverConn, clientConn)
}

func decrypt(conn net.Conn) error {
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
			clientConn.Close()
			return err
		}

		fmt.Println("tls decrypted")
		return nil
	}

	err := fmt.Errorf("error in handshake")
	clientConn.Close()
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

func tunnel(from, to io.ReadWriteCloser) {
	defer from.Close()
	defer to.Close()
	io.Copy(from, to)
	fmt.Println("copied data from", from, "to", to)
}