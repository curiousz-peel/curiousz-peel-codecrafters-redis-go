package main

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

type respHandler struct {
	data     []byte
	commands map[string]string
	conn     net.Conn
}

func InitRESP(conn net.Conn) *respHandler {
	return &respHandler{conn: conn, commands: make(map[string]string)}
}

func (r *respHandler) Read() error {
	buf := make([]byte, 128)
	_, err := r.conn.Read(buf)
	if err != nil {
		fmt.Printf("couldn't read from %v\n", r.conn.LocalAddr().String())
		return err
	}
	r.data = buf
	return nil
}

func (r *respHandler) Parse() {
	var isArg bool
	var currentCommand string
	splitRe := regexp.MustCompile(`\$(\d+)\r\n(\w+)`)
	trimRe := regexp.MustCompile(`\$(\d+)\r\n`)

	tokens := splitRe.FindAll(r.data, -1)
	if tokens != nil {
		for _, token := range tokens {
			token = trimRe.ReplaceAll(token, []byte(""))
			if !isArg {
				currentCommand = string(token)
				r.commands[currentCommand] = ""
			} else {
				r.commands[currentCommand] = string(token)
			}
			isArg = true
		}
		return
	}
	fmt.Printf("no tokens after regex expression: %s\n", splitRe.String())
}

func (r respHandler) handleCommand(command string, argument string) (response []byte, err error) {
	switch strings.ToLower(command) {
	case "ping":
		response = []byte("+PONG\r\n")
	case "echo":
		response = []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(argument), argument))
	default:
		err = fmt.Errorf("unknown command %s", command)
	}
	return
}

func (r respHandler) Execute() error {
	if r.commands != nil {
		for comm, arg := range r.commands {
			resp, err := r.handleCommand(comm, arg)
			if err != nil {
				return err
			}
			_, err = r.conn.Write(resp)
			if err != nil {
				return fmt.Errorf("couldn't write to %s: %v", r.conn.LocalAddr().String(), err.Error())
			}
		}
		return nil
	}
	return fmt.Errorf("no commands to execute")
}
