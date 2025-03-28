//go:build host
// +build host

package controllers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"server2/api/v1/endpoints"
	"server2/application/adapters"
	"server2/application/entities"
	usecases "server2/application/useCases"
	"server2/application/useCases/commands"
	"server2/application/useCases/handlers"
	"server2/security"

	"github.com/gorilla/websocket"
)

const IsHost = true

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (Change this for security)
	},
}

var nodeRepo adapters.InMemoryNodeRepository = adapters.NewInMemoryNodeRepository()
var responseRepo adapters.InMemoryResponseRepository = adapters.NewInMemoryResponseRepository()

var saveNodeUseCase = usecases.SaveNodeUseCase{Repo: &nodeRepo}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		nodeID := conn.RemoteAddr().String()
		node := nodeRepo.Get(nodeID)

		parseCommand := usecases.ParseCommandUseCase{Cipher: &node.Cipher}
		command, err := parseCommand.Execute(string(msg))

		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("[received]", command)
		if command.Entry == "connect" {
			fmt.Println("connecting")
			sharedKey, _, _ := commands.Connect(command.Id, command.Args)
			cipher := security.CreateCipher(sharedKey)
			saveNodeUseCase.Execute(entities.Node{Conn: conn, Name: "any", Cipher: cipher})
			encodedPublicKey := base64.RawStdEncoding.EncodeToString(security.PublicKey.Bytes())
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("_ connect %s %s", encodedPublicKey, "host")))
			continue
		}

		handleCommands := handlers.HandleCommandsUseCase{ResponseRepo: &responseRepo}
		_, err = handleCommands.Execute(command)

		if err != nil {
			fmt.Println("command error:", err)
			continue
		}
	}
}

func RunWebsocketAsHost() {

	v1Handlers := endpoints.V1Handlers{NodeRepo: &nodeRepo, ResponseRepo: &responseRepo}

	http.HandleFunc("/ws", handleWebSocket)

	http.HandleFunc("/auth/login", endpoints.AuthLoginHandler)

	http.HandleFunc("/auth/token", endpoints.AuthClientToken)

	http.HandleFunc("/auth/register", endpoints.AuthRegisterHandler)

	http.HandleFunc("/auth/status", endpoints.AuthHasUserHandler)

	http.HandleFunc("/auth/reset", endpoints.AuthResetAccountHandler)

	http.HandleFunc("/v1/connections", v1Handlers.ConnectionsHandler)

	http.HandleFunc("/v1/connections/{connection}/blocks", v1Handlers.BlockAddressHandler)

	http.HandleFunc("/v1/connections/{connection}/redirects", v1Handlers.RedirectAddressHandler)

	http.HandleFunc("/v1/connections/{connection}/reload", v1Handlers.ReloadHandler)

	http.HandleFunc("/v1/connections/{connection}/confighash", v1Handlers.ConfigHashHandler)

	fmt.Println("WebSocket server running on ws://localhost:8080/ws")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

func RunWebsocketAsNode() {
	panic("not implemented")
}
