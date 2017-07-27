package main

import (
	"net"
	"os"
	"fmt"
	"strconv"
	"log"
	"errors"
)

func main() {
	port,err := port(os.Args)
	if err != nil {
		log.Fatal("Please specify a port to listen")
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":" + port)
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

func port(args []string) (string, error) {
	if len(args) > 1 && args[1] != "" {
		port := args[1]
		_, err := strconv.ParseUint(port, 10, 0)
		if err != nil {
			log.Fatal("Specify a port to listen")
		}
		return port, nil
	}
	return "", errors.New("value out of range")
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
		println("[x]: " + string(buf[0:n]))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
