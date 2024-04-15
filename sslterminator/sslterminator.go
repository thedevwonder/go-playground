package sslterminator

import (
	"fmt"
	"crypto/tls"
	"net"
	"io"
	"runtime"
	"time"
	"os/exec"
)

func Helper () {

	runtime.GOMAXPROCS(1)

	certPath := "/Users/theboywonder/Documents/projects/go-playground/sslterminator/certs/ecdsa_cert.pem"
	keyPath := "/Users/theboywonder/Documents/projects/go-playground/sslterminator/certs/ecdsa_key.pem"

	fmt.Println("Starting ssl terminator")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err!=nil {
		fmt.Println("error in loading cert", err)
		return
	}

	tlsConfig := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	ls, err := tls.Listen("tcp", ":443", &tlsConfig)

	if err != nil {
		fmt.Println("error in creating tls server", err)
		return
	}
	fmt.Println("client is running on port: 443", ls.Addr())

	go func() {
		for {
			t1 := time.Now()
			conn, err := ls.Accept()
			t2 := time.Now()
			fmt.Println(t2.Sub(t1), "time taken for accept")
			if err != nil {
				fmt.Println("error in accepting client connections")
				//resource released
				conn.Close()
				return
			}
			go sslTerminator(conn)
		}
	}()

	arg0 := "openssl"
	arg1 := "s_client"
	arg2 := "localhost:443"
	cmd := exec.Command(arg0, arg1, arg2)
	_, sslErr := cmd.Output()
	if sslErr != nil {
		fmt.Println("error in ssl client", sslErr)
	}
}

func AcceptConnections(ls net.Listener) {
	for {
		conn, err := ls.Accept()
		if err != nil {
			fmt.Println("error in accepting client connections")
			return
		}

		Handshake(conn)
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