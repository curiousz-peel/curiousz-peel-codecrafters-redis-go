package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			// continue
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 128)
	for {
		bitesRead, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("couldn't read from %v\n", conn.LocalAddr().String())
			return
		}

		fmt.Println("received data", buf[:bitesRead])

		resp := []byte("+PONG\r\n")
		_, err = conn.Write(resp)
		if err != nil {
			fmt.Printf("couldn't write to %v\n", conn.LocalAddr().String())
		}
	}

}
