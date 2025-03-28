//go:build !host
// +build !host

package controllers

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"server2/application/adapters"
	"server2/application/entities"
	usecases "server2/application/useCases"
	"server2/application/useCases/commands"
	"server2/application/useCases/handlers"
	"server2/enviroment"
	"server2/security"
	"server2/utils"

	"github.com/gorilla/websocket"
)

func connectWebsocket() *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", enviroment.MAIN_SERVER_ADDRESS), nil)
	if err != nil {
		fmt.Println("Connection error:", err)
		return nil
	}
	return conn
}

// var nodeRepo adapters.InMemoryNodeRepository = adapters.NewInMemoryNodeRepository()
var responseRepo adapters.InMemoryResponseRepository = adapters.NewInMemoryResponseRepository()

// var saveNode = usecases.SaveNodeUseCase{Repo: &nodeRepo}

const IsHost = false

func RunWebsocketAsNode() {
	name := utils.GetEnvOrDefault("NAME", fmt.Sprintf("%x", rand.Int()))
	var conn *websocket.Conn
	fmt.Println("trying connection")
	conn = connectWebsocket()
	host := entities.Node{Conn: conn, Name: "Host", Cipher: nil}

	parse := usecases.ParseCommandUseCase{Cipher: &host.Cipher}

	HandleCommands := handlers.HandleCommandsUseCase{ResponseRepo: &responseRepo}

	go func() {
		fmt.Println("sending and receiving keys")

		var encodedPublicKey = base64.RawStdEncoding.EncodeToString(security.PublicKey.Bytes())

		responseId := fmt.Sprintf("%x", rand.Int())
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s connect %s %s", responseId, encodedPublicKey, name)))
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			conn.Close()
			conn = nil
			break
		}

		command, err := parse.Execute(string(msg))
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("[received]", command.IsEncrypted, command.String())

		if command.Entry == "connect" {
			fmt.Println("connecting")
			sharedKey, _, _ := commands.Connect(command.Id, command.Args)
			cipher := security.CreateCipher(sharedKey)
			host.Cipher = cipher
			continue
		}

		response, err := HandleCommands.Execute(command)

		fmt.Println("response=", response, "\nerror=", err)
		if err != nil {
			fmt.Println("error")
			continue
		}
		host.Send("_ add response "+command.Id+" "+response, true)

	}
}

func RunWebsocketAsHost() {
	panic("not implemented")
}
