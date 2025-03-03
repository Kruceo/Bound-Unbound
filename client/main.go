package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unbound-side-client/commands"
	"unbound-side-client/host"
	"unbound-side-client/utils"

	"github.com/gorilla/websocket"
)

func connectWebsocket() *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", host.MAIN_SERVER_ADDRESS), nil)
	if err != nil {
		fmt.Println("Connection error:", err)
		return nil
	}
	return conn
}

func main() {
	name := utils.GetEnvOrDefault("NAME", fmt.Sprintf("%x", rand.Int()))
	// Connect to WebSocket server
	conn := connectWebsocket()
	for conn == nil {
		conn = connectWebsocket()
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Connected to WebSocket server")
	conn.WriteMessage(websocket.TextMessage, []byte("add connection "+name))

	scanner := bufio.NewScanner(os.Stdin)

	// Read messages in a separate goroutine
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Read error:", err)
				conn.Close()
				conn = nil
				for conn == nil {
					conn = connectWebsocket()
					time.Sleep(1 * time.Second)
				}
				conn.WriteMessage(websocket.TextMessage, []byte("add connection "+name))
				// fmt.Println("Connection")
			} else {
				HandleCommands(conn, string(msg))
			}
		}
	}()

	// Read input from user and send to WebSocket

	for {
		fmt.Print("\n>: ")
		scanner.Scan()
		text := scanner.Text()
		HandleCommands(conn, "local "+text)
	}
}

func HandleCommands(conn *websocket.Conn, str string) {
	fmt.Println("[command]", str)
	splt := strings.Split(str, " ")
	if len(splt) < 2 {
		fmt.Println("Wrong Syntax: (" + str + ")")
		return
	}
	id := splt[0]
	command := splt[1]
	args := splt[2:]
	if command == "block" {
		err := commands.Block(conn, id, false, args)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if command == "unblock" {
		err := commands.Block(conn, id, true, args)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if command == "list" {
		err := commands.List(conn, id, args)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if command == "reload" {
		err := commands.ReloadConfig(conn, id)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
