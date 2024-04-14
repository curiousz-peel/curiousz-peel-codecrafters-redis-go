package main

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// type cacheVal struct {
// 	value  string
// 	expiry time.Time
// }

type respHandler struct {
	data     []byte
	commands map[string][]string
	cache    map[string]string
	conn     net.Conn
}

func InitRESP(conn net.Conn) *respHandler {
	return &respHandler{conn: conn, commands: make(map[string][]string), cache: make(map[string]cacheVal)}
}

func (r *respHandler) Read() error {
	buf := make([]byte, 128)
	_, err := r.conn.Read(buf)
	if err != nil {
		return err
	}
	r.data = buf
	return nil
}

func (r *respHandler) Parse() {
	splitRe := regexp.MustCompile(`\r\n`)

	tokens := splitRe.Split(string(r.data), -1)
	if tokens != nil {
		wordCount, err := strconv.Atoi(tokens[0][1:])
		if err != nil {
			fmt.Println("couldn't parse length of words")
			return
		}
		r.commands[tokens[2]] = []string{}
		for i := 2; i <= wordCount; i++ {
			r.commands[tokens[2]] = append(r.commands[tokens[2]], tokens[i*2])
		}
		return
	}
	fmt.Printf("no tokens after regex expression: %s\n", splitRe.String())
}

func (r respHandler) handleCommand(command string, arguments []string) (response []byte, err error) {
	switch strings.ToLower(command) {
	case "ping":
		response = []byte("+PONG\r\n")
	case "echo":
		response = []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(arguments[0]), arguments[0]))
	case "set":
		fmt.Println(arguments)
		r.cache[arguments[0]] = arguments[1]
		response = []byte("+OK\r\n")
	case "get":
		val, ok := r.cache[arguments[0]]
		if !ok {
			response = []byte("$-1\r\n")
		} else {
			response = []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val))
		}
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
			defer delete(r.commands, comm)
			_, err = r.conn.Write(resp)
			if err != nil {
				return fmt.Errorf("couldn't write to %s: %v", r.conn.LocalAddr().String(), err.Error())
			}
		}
		return nil
	}
	return fmt.Errorf("no commands to execute")
}
