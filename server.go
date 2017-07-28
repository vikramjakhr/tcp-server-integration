package main

import (
	"net"
	"os"
	"fmt"
	"strconv"
	"log"
	"github.com/vikramjakhr/tcp-server/util"
	"github.com/vikramjakhr/tcp-server/mqtt"
)

const defaultPort = "6666"

func main() {
	port := listeningPort(os.Args)
	registerTCPListener(port)
}

func listeningPort(args []string) (string) {
	if len(args) > 1 && args[1] != "" {
		port := args[1]
		_, err := strconv.ParseUint(port, 10, 0)
		if err != nil {
			log.Fatal("Specify a port to listen")
		}
		return port
	}
	log.Println("No port specified, Falling back to default ", defaultPort)
	return defaultPort
}

func registerTCPListener(port string) {
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
		handleClient(conn)
		conn.Close()
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
