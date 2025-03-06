package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var Connections map[string]*websocket.Conn = make(map[string]*websocket.Conn)
var Responses map[string]string = make(map[string]string)
var ResponseCH = make(chan string)

func WaitForResponse(id string) {
	for {
		select {
		case t := <-ResponseCH:
			if t == id {
				return
			}
		}
	}
}

func Add(conn *websocket.Conn, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("wrong Syntax: (%v)\nuse add <response|connection> data", args)
	}
	if args[0] == "connection" {
		name := strings.Join(args[1:], " ")
		// add connection(arg0) <id>(arg1)
		Connections[name] = conn
		fmt.Println(name, "connected")

		go func() {
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				err := conn.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					delete(Connections, name)
					fmt.Println(name, "disconnected")
					break
				}
			}
		}()

	}

	if args[0] == "response" {
		// add response(arg0) <id>(arg1) data(rest args concatened)
		Responses[args[1]] = strings.Join(args[2:], " ")
		ResponseCH <- args[1]
	}

	// fmt.Println(connections)
	return nil
}
