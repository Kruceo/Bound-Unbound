package memory

import (
	"crypto/cipher"
	"fmt"
	"time"
	"unbound-mngr-host/security"

	"github.com/gorilla/websocket"
)

func GetHost() (Client, error) {
	if client, exists := Connections["_HOST_"]; exists {
		return client, nil
	}
	return Client{}, fmt.Errorf("no host registered")
}
func SetHost(conn *websocket.Conn, cipher cipher.AEAD) error {
	if _, exists := Connections["_HOST_"]; !exists {
		Connections["_HOST_"] = Client{Conn: conn, Cipher: cipher}
		return nil
	}
	return fmt.Errorf("host already registered")
}

type Client struct {
	Conn   *websocket.Conn
	Name   string
	Cipher cipher.AEAD
}

func (client Client) Send(msg string, encrypt bool) {
	if encrypt {
		base64Msg := security.CipherMessageBase64(msg, client.Cipher)
		client.Conn.WriteMessage(websocket.TextMessage, append([]byte("#$"), base64Msg...))
	} else {
		client.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

var Connections map[string]Client = make(map[string]Client)
var Responses map[string]string = make(map[string]string)
var ResponseCH = make(chan string)

func ReadResponse(id string) string {
	response := Responses[id]
	delete(Responses, id)
	return response
}

func WaitForResponse(id string) error {
	go func() {
		ticker := time.NewTimer(30 * time.Second)
		defer ticker.Stop()

		<-ticker.C
		fmt.Println("timeout for", id)
		ResponseCH <- "_TIMEOUT_"
	}()

	if _, exists := Responses[id]; exists {
		return nil
	}

	for t := range ResponseCH {
		fmt.Println("t=", t)
		if t == id {
			break
		}
		if t == "_TIMEOUT_" {
			return fmt.Errorf("timeout for response id %s", id)
		}
	}
	return nil
}
