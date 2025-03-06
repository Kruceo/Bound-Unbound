package main

import (
	"fmt"
	"net/http"
	"strings"
	"unbound-mngr-host/api/v1/handlers"
	"unbound-mngr-host/commands"

	"github.com/gorilla/websocket"
)

// Upgrader to upgrade HTTP to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (Change this for security)
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected")

	for {
		// Read message from client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		fmt.Println("Received:", string(msg))
		splt := strings.Split(string(msg), " ")
		if len(splt) < 1 {
			fmt.Println("wrong syntax: " + string(msg))
			continue
		}
		command := splt[0]
		args := splt[1:]

		if command == "add" {
			commands.Add(conn, args)
		}

		// err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	http.HandleFunc("/v1/connections", handlers.ConnectionsHandler)

	http.HandleFunc("/v1/connections/{connection}/blocked", handlers.BlockAddressHandler)

	http.HandleFunc("/v1/connections/{connection}/redirects", handlers.RedirectAddressHandler)

	http.HandleFunc("/v1/connections/{connection}/reload", handlers.ReloadHandler)

	http.HandleFunc("/v1/connections/{connection}/confighash", handlers.ConfigHashHandler)

	fmt.Println("WebSocket server running on ws://localhost:8080/ws")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
