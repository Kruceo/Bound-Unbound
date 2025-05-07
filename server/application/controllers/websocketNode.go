//go:build !host
// +build !host

package controllers

import (
	"crypto/cipher"
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
	"server2/application/infrastructure"
	usecases "server2/application/useCases"

	"server2/application/useCases/handlers"
	"server2/application/useCases/security"

	"github.com/gorilla/websocket"
)

const IsHost = false

type WebsocketClientController struct {
	name                 string
	nodeID               *usecases.NodeIDUseCase
	responseRepo         infrastructure.ResponsesReporisory
	cipherCreation       security.CiphersUseCase
	publicKey            *ecdh.PublicKey
	cipherCommandMessage usecases.CipherCommandMessageUseCase
	hostConn             *websocket.Conn
	handleCommands       handlers.HandleCommandsUseCase
	mainCipher           *cipher.AEAD
	sharedKeyCreation    security.CreateSharedKeyUseCase
}

func NewWebsocketClientController(name string, hostConn *websocket.Conn, responseRepo infrastructure.ResponsesReporisory, privateKey *ecdh.PrivateKey, publicKey *ecdh.PublicKey) WebsocketClientController {
	handleCommands := handlers.HandleCommandsUseCase{ResponseRepo: responseRepo}
	cuc := security.CiphersUseCase{}
	cmuc := usecases.NewCipherMessageUseCase()
	skuc := security.NewCreateSharedKeyUseCase(*privateKey)
	return WebsocketClientController{
		name:                 name,
		nodeID:               usecases.NewNodeIDUseCase(),
		responseRepo:         responseRepo,
		handleCommands:       handleCommands,
		cipherCreation:       cuc,
		cipherCommandMessage: cmuc,
		publicKey:            publicKey,
		sharedKeyCreation:    skuc,
		hostConn:             hostConn,
	}
}

func (wsc *WebsocketClientController) ExecuteStringAsCommand(cmdStr string) error {
	parse := usecases.ParseCommandUseCase{Cipher: wsc.mainCipher}
	command, err := parse.Execute(cmdStr)
	if err != nil {
		return err
	}

	fmt.Printf("[received %v] %s\n", command.IsEncrypted, command.String())

	if command.Entry == "connect" && len(command.Args) >= 2 {
		sharedKey, err := wsc.sharedKeyCreation.Execute(command.Args[0])
		if err != nil {
			return err
		}
		newCipher := wsc.cipherCreation.CreateCipher(sharedKey)
		wsc.mainCipher = &newCipher
		return nil
	}

	response, err := wsc.handleCommands.Execute(command)
	if err != nil {
		return err
	}
	return wsc.SendEncryptedResponse(command.Id, response)
}

func (wsc *WebsocketClientController) Connect() error {
	fmt.Println("connecting with host")
	thisNodeID, err := wsc.nodeID.ReadOrCreateFile()

	if err != nil {
		return err
	}

	var encodedPublicKey = base64.RawStdEncoding.EncodeToString(wsc.publicKey.Bytes())
	err = wsc.hostConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("_ connect %s %s %s", encodedPublicKey, thisNodeID, wsc.name)))
	return err

}

func (wsc *WebsocketClientController) SetConnection(conn *websocket.Conn) {
	wsc.hostConn = conn
}

func (wsc *WebsocketClientController) HasConnection() bool {
	return wsc.hostConn != nil
}

func (wsc *WebsocketClientController) ReadConn() (string, error) {
	_, content, err := wsc.hostConn.ReadMessage()
	return string(content), err
}

func (wsc *WebsocketClientController) SendEncryptedResponse(id string, str string) error {
	encryptedMessage, err := wsc.cipherCommandMessage.Execute(fmt.Sprintf("_ add response %s %s", id, str), wsc.mainCipher)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return err
	}

	wsc.hostConn.WriteMessage(websocket.TextMessage, encryptedMessage)
	return nil
}
