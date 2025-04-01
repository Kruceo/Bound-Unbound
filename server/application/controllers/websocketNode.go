//go:build !host
// +build !host

package controllers

import (
	"crypto/cipher"
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
	"math/rand"
	"server2/application/adapters"
	usecases "server2/application/useCases"

	"server2/application/useCases/handlers"
	"server2/application/useCases/security"
	"server2/enviroment"

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
// var getOrCreate = usecases.CreateNodeUseCase{}
var cipherCreation = security.CiphersUseCase{}

const IsHost = false

var publicKey *ecdh.PublicKey
var createSharedKey security.CreateSharedKeyUseCase

func init() {
	genKeysUseCase := security.GenKeysUseCase{}
	priv, pub := genKeysUseCase.GenKeys()
	publicKey = pub
	createSharedKey = security.NewCreateSharedKeyUseCase(*priv)
}

func RunWebsocketAsNode() {
	name := utils.GetEnvOrDefault("NAME", fmt.Sprintf("%x", rand.Int()))

	var conn *websocket.Conn
	var connLocker sync.Mutex = sync.Mutex{}
	var cipher cipher.AEAD
	parse := usecases.ParseCommandUseCase{Cipher: &cipher}
	cipherMessage := usecases.CipherMessageUseCase{}
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

				var encodedPublicKey = base64.RawStdEncoding.EncodeToString(publicKey.Bytes())

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

		if command.Entry == "connect" && len(command.Args) >= 2 {
			fmt.Println("connecting")
			sharedKey, err := createSharedKey.Execute(command.Args[0])
			if err != nil {
				panic(err)
			}
			cipher = cipherCreation.CreateCipher(sharedKey)
			continue
		}

		response, err := HandleCommands.Execute(command)

		fmt.Println("response=", response, "\nerror=", err)
		if err != nil {
			fmt.Println("error")
			continue
		}
		connLocker.Lock()
		encryptedMessage, err := cipherMessage.Execute(fmt.Sprintf("_ add response %s %s", command.Id, response), &cipher)
		if err != nil {
			connLocker.Unlock()
			fmt.Println(err)
			continue
		}
		conn.WriteMessage(websocket.TextMessage, encryptedMessage)
		connLocker.Unlock()
		fmt.Println(err)
	}
}

func RunWebsocketAsHost() {
	panic("not implemented")
}
