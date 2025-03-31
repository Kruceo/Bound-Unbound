//go:build host
// +build host

package controllers

import (
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"server2/application/adapters"
	"server2/application/controllers/api/v1/endpoints"

	usecases "server2/application/useCases"
	"server2/application/useCases/commands"
	"server2/application/useCases/handlers"
	"server2/application/useCases/security"

	"time"

	"github.com/gorilla/websocket"
)

var privateKey *ecdh.PrivateKey
var publicKey *ecdh.PublicKey

func init() {
	genKeysUseCase := security.GenKeysUseCase{}
	priv, pub := genKeysUseCase.GenKeys()
	privateKey = priv
	publicKey = pub
}

const IsHost = true

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (Change this for security)
	},
}

var nodeRepo adapters.InMemoryNodeRepository = adapters.NewInMemoryNodeRepository()
var responseRepo adapters.InMemoryResponseRepository = adapters.NewInMemoryResponseRepository()

var saveNodeUseCase = usecases.SaveNodeUseCase{Repo: &nodeRepo}
var deleteNodeUseCAse = usecases.DeleteNodeUseCase{Repo: &nodeRepo}
var getOrCreate = usecases.GetOrCreateUseCase{Repo: &nodeRepo}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	cipherCreation := security.CiphersUseCase{}
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
		node, err := getOrCreate.Execute(nodeID, conn)
		if err != nil {
			fmt.Println("Node error:", err)
			break
		}
		parseCommand := usecases.ParseCommandUseCase{Cipher: &node.Cipher}
		command, err := parseCommand.Execute(string(msg))

		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("[received]", command)
		if command.Entry == "connect" && len(command.Args) >= 2 {
			fmt.Println("connecting")
			sharedKey, name, _ := commands.Connect(
				privateKey,
				command.Id,
				strings.Join(command.Args[1:], " "),
				command.Args[0],
			)
			cipher := cipherCreation.CreateCipher(sharedKey)
			saveNodeUseCase.Execute(conn, name, cipher)
			encodedPublicKey := base64.RawStdEncoding.EncodeToString(publicKey.Bytes())
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("_ connect %s %s", encodedPublicKey, "host")))

			go func() {
				ticker := time.NewTicker(5 * time.Second)
				defer ticker.Stop()
				for range ticker.C {
					err := conn.WriteMessage(websocket.PingMessage, nil)
					if err != nil {

						if deleteNodeUseCAse.Execute(conn.RemoteAddr().String()) != nil {
							fmt.Println("problem to remove node from repository")
							break
						}
						fmt.Println(name, "disconnected")
						break
					}
				}
			}()

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
