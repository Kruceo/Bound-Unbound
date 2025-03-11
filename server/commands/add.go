package commands

import (
	"fmt"
	"strings"
	"unbound-mngr-host/memory"

	"github.com/gorilla/websocket"
)

func Add(conn *websocket.Conn, id string, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("wrong Syntax: (%v)\nuse 'id' add <response|connection> data", args)
	}

	if args[0] == "response" {
		// add response(arg0) <id>(arg1) data(rest args concatened)
		memory.Responses[args[1]] = strings.Join(args[2:], " ")
		memory.ResponseCH <- args[1]
	}

	// fmt.Println(connections)
	return nil
}
