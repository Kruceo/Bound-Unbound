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
	"sync"
	"time"

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
	var connLocker sync.Mutex = sync.Mutex{}

	host := entities.Node{Conn: conn, Name: "Host", Cipher: nil}

	parse := usecases.ParseCommandUseCase{Cipher: &host.Cipher}

	HandleCommands := handlers.HandleCommandsUseCase{ResponseRepo: &responseRepo}

	for {
		if conn == nil {
			fmt.Println("trying connection")
			conn = connectWebsocket()
			go func() {
				if conn == nil {
					return
				}
				fmt.Println("sending and receiving keys")

				var encodedPublicKey = base64.RawStdEncoding.EncodeToString(security.PublicKey.Bytes())

				responseId := fmt.Sprintf("%x", rand.Int())
				connLocker.Lock()
				defer connLocker.Unlock()
				conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s connect %s %s", responseId, encodedPublicKey, name)))
				// host.Send(, false)
			}()
			if conn == nil {
				time.Sleep(3 * time.Second)
			}
			continue
		}
		fmt.Println("listening commands")
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			conn.Close()
			conn = nil
			continue
		}

		command, err := parse.Execute(string(msg))
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("[received %v] %s\n", command.IsEncrypted, command.String())

		if command.Entry == "connect" {
			fmt.Println("connecting")
			sharedKey, _, _ := commands.Connect(command.Id, command.Args)
			cipher := security.CreateCipher(sharedKey)
			host.Cipher = cipher
			host.Conn = conn
			continue
		}

		response, err := HandleCommands.Execute(command)

		fmt.Println("response=", response, "\nerror=", err)
		if err != nil {
			fmt.Println("error")
			continue
		}
		connLocker.Lock()
		err = host.Send("_ add response "+command.Id+" "+response, true)
		connLocker.Unlock()
		fmt.Println(err)
	}
}

func RunWebsocketAsHost() {
	panic("not implemented")
}
