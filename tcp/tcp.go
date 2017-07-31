package tcp

import (
	"net"
	"fmt"
	"os"
	"gitlab.intelligrape.net/tothenew/tcp-server-integration/mqtt"
	"gitlab.intelligrape.net/tothenew/tcp-server-integration/util"
)

func RegisterTCPListener(port string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":"+port)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	println("Server listening to tcp port: " + port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go func() {
			defer conn.Close()
			handleClient(conn)
		}()
	}
}

func handleClient(conn net.Conn) {
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}

		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			return
		}
		data := string(buf[0:n])
		println("[x]: " + data)
		mqtt.Publish(util.ParsePayload(data))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
