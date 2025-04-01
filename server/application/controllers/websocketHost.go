//go:build host
// +build host

package controllers

import (
	"crypto/cipher"
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"server2/application/entities"

	usecases "server2/application/useCases"

	"server2/application/useCases/handlers"
	"server2/application/useCases/security"

	"github.com/gorilla/websocket"
)

const IsHost = true

type HostController struct {
	upgrader             websocket.Upgrader
	nodeRepo             entities.NodeRepository
	responseRepo         entities.ResponsesReporisory
	saveNode             usecases.SaveNodeUseCase
	deleteNode           usecases.DeleteNodeUseCase
	getOrCreate          usecases.GetOrCreateUseCase
	getNode              usecases.GetNodeUseCase
	publicKey            ecdh.PublicKey
	mainCipher           *cipher.AEAD
	cipherCommandMessage usecases.CipherCommandMessageUseCase
	sharedKeyCreation    security.CreateSharedKeyUseCase
	ciphersCreation      security.CiphersUseCase
	handleCommands       handlers.HandleCommandsUseCase
}

func NewHostController(nodeRepo entities.NodeRepository, responseRepo entities.ResponsesReporisory, privateKey ecdh.PrivateKey, publicKey ecdh.PublicKey) HostController {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all connections (Change this for security)
		},
	}

	var saveNodeUseCase = usecases.SaveNodeUseCase{Repo: &nodeRepo}
	var deleteNodeUseCase = usecases.DeleteNodeUseCase{Repo: &nodeRepo}
	var getOrCreateNode = usecases.GetOrCreateUseCase{Repo: &nodeRepo}
	var getNode = usecases.GetNodeUseCase{Repo: &nodeRepo}
	skuc := security.NewCreateSharedKeyUseCase(privateKey)
	cuc := security.CiphersUseCase{}
	var commandHandler = handlers.HandleCommandsUseCase{ResponseRepo: responseRepo}
	return HostController{
		upgrader:             upgrader,
		saveNode:             saveNodeUseCase,
		deleteNode:           deleteNodeUseCase,
		nodeRepo:             nodeRepo,
		responseRepo:         responseRepo,
		getNode:              getNode,
		getOrCreate:          getOrCreateNode,
		publicKey:            publicKey,
		mainCipher:           nil,
		cipherCommandMessage: usecases.NewCipherMessageUseCase(),
		sharedKeyCreation:    skuc,
		ciphersCreation:      cuc,
		handleCommands:       commandHandler,
	}

}

func (wsc *HostController) SendEncryptedMessageToNode(nodeId string, id string, str string) error {
	encryptedMessage, err := wsc.cipherCommandMessage.Execute(fmt.Sprintf("_ add response %s %s", id, str), wsc.mainCipher)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return err
	}

	node := wsc.getNode.Execute(nodeId)
	if node == nil {
		return fmt.Errorf("node not found: %s", nodeId)
	}

	return node.Conn.WriteMessage(websocket.TextMessage, encryptedMessage)
}

func (wsc *HostController) ExecuteStringAsCommand(cmdStr string, conn *websocket.Conn, cipher *cipher.AEAD) error {
	parseCommand := usecases.ParseCommandUseCase{Cipher: cipher}
	command, err := parseCommand.Execute(string(cmdStr))
	if err != nil {
		fmt.Println(err)
		return err
	}

	// fmt.Printf("[received %v] %s\n", command.IsEncrypted, command.String())
	fmt.Println(command.Entry)
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++=")
	if command.Entry == "connect" && len(command.Args) >= 2 {
		fmt.Println("receiving connection")
		name := strings.Join(command.Args[1:], " ")
		sharedKey, err := wsc.sharedKeyCreation.Execute(command.Args[0])
		if err != nil {
			return err
		}
		cipher := wsc.ciphersCreation.CreateCipher(sharedKey)
		nodeID, err := wsc.saveNode.Execute(conn, name, cipher)
		if err != nil {
			return err
		}
		wsc.SendConnectToNode(nodeID)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = wsc.handleCommands.Execute(command)
	if err != nil {
		return err
	}
	return nil
}

func (wsc *HostController) SendConnectToNode(nodeId string) error {
	fmt.Println("connecting with", nodeId)
	var encodedPublicKey = base64.RawStdEncoding.EncodeToString(wsc.publicKey.Bytes())

	node := wsc.getNode.Execute(nodeId)
	if node == nil {
		return fmt.Errorf("node not found: %s", nodeId)
	}

	err := node.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("_ connect %s %s", encodedPublicKey, "host")))
	return err
}

func (wsc *HostController) AddNodeToRepo(conn *websocket.Conn, name string, cipher cipher.AEAD) (string, error) {
	return wsc.saveNode.Execute(conn, name, cipher)
}

func (wsc *HostController) OnMessageHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsc.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		nodeID := conn.RemoteAddr().String()
		node, err := wsc.getOrCreate.Execute(nodeID, conn)
		if err != nil {
			fmt.Println("Node error:", err)
			break
		}

		if err := wsc.ExecuteStringAsCommand(string(msg), conn, &node.Cipher); err != nil {
			fmt.Println("Command error:", err)
			break
		}
	}
}
