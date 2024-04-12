package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	respHandler := InitRESP(conn)
	for {
		err := respHandler.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			fmt.Println(err)
			os.Exit(2)

		}
		respHandler.Parse()
		err = respHandler.Execute()
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
	}
}
