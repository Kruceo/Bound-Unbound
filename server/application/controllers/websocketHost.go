package controllers

import (
	"crypto/cipher"
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"server2/application/infrastructure"

	usecases "server2/application/useCases"

	"server2/application/useCases/handlers"
	"server2/application/useCases/security"

	"github.com/gorilla/websocket"
)

type HostController struct {
	upgrader             websocket.Upgrader
	responseRepo         infrastructure.ResponsesReporisory
	nodePersistence      *usecases.NodePersistenceUseCase
	publicKey            ecdh.PublicKey
	mainCipher           *cipher.AEAD
	cipherCommandMessage usecases.CipherCommandMessageUseCase
	sharedKeyCreation    security.CreateSharedKeyUseCase
	ciphersCreation      security.CiphersUseCase
	handleCommands       handlers.HandleCommandsUseCase
}

func NewHostController(nodeRepo infrastructure.NodeRepository, responseRepo infrastructure.ResponsesReporisory, privateKey ecdh.PrivateKey, publicKey ecdh.PublicKey) HostController {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all connections (Change this for security)
		},
	}

	nodePersistence := usecases.NewNodePersistenceUseCase(nodeRepo)

	skuc := security.NewCreateSharedKeyUseCase(privateKey)
	cuc := security.CiphersUseCase{}
	var commandHandler = handlers.HandleCommandsUseCase{ResponseRepo: responseRepo}
	return HostController{
		nodePersistence:      nodePersistence,
		upgrader:             upgrader,
		responseRepo:         responseRepo,
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

	node, err := wsc.nodePersistence.Get(nodeId)
	if err != nil {
		return fmt.Errorf("node not found: %s", nodeId)
	}

	return node.Conn.WriteMessage(websocket.TextMessage, encryptedMessage)
}

func (wsc *HostController) ExecuteStringAsCommand(cmdStr string, conn *websocket.Conn) error {
	// use remote address as "nodeID" because is much more
	// easy gets remote address at each websocket call than
	// include a "nodeid" header at each call;
	// in middleware we can get the real
	// nodeID and compare as we need
	remoteAddress := conn.RemoteAddr().String()
	var messageCipher *cipher.AEAD
	node, err := wsc.nodePersistence.FindOneByRemoteAddress(remoteAddress)
	if err != nil {
		fmt.Println(err)
	} else {
		messageCipher = node.Cipher
	}

	parseCommand := usecases.ParseCommandUseCase{Cipher: messageCipher}
	command, err := parseCommand.Execute(string(cmdStr))
	if err != nil {
		fmt.Println(err)
		return err
	}

	// fmt.Printf("[received %v] %s\n", command.IsEncrypted, command.String())

	if command.Entry == "connect" && len(command.Args) >= 3 {
		fmt.Println("receiving connection from", remoteAddress)
		name := strings.Join(command.Args[2:], " ")
		nodeID := command.Args[1]
		sharedKey, err := wsc.sharedKeyCreation.Execute(command.Args[0])
		if err != nil {
			return err
		}
		cipher := wsc.ciphersCreation.CreateCipher(sharedKey)

		wsc.nodePersistence.Save(nodeID, name, conn, &cipher)
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

	node, err := wsc.nodePersistence.Get(nodeId)
	if err != nil {
		return fmt.Errorf("node not found: %s", nodeId)
	}

	err = node.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("_ connect %s %s", encodedPublicKey, "host")))
	return err
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
			findedNode, err := wsc.nodePersistence.FindOneByRemoteAddress(conn.RemoteAddr().String())
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = wsc.nodePersistence.Delete(findedNode.ID)
			if err != nil {
				fmt.Println("delete error:", err)
				continue
			}
			break
		}

		if err := wsc.ExecuteStringAsCommand(string(msg), conn); err != nil {
			fmt.Println("Command error:", err)
			break
		}
	}
}
