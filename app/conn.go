package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

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
