package memory

import (
	"crypto/cipher"
	"encoding/base64"
	"fmt"
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
		nonce := security.RandomNonce()
		encryptedContent := client.Cipher.Seal(nil, nonce, []byte(msg), nil)
		formatedMsg := encryptedContent
		formatedMsg = append(formatedMsg, nonce...)
		// fmt.Println("DECODED", formatedMsg)
		base64Msg := base64.RawStdEncoding.EncodeToString(formatedMsg)
		// fmt.Println("ENCODED", []byte(base64Msg))
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

func WaitForResponse(id string) {
	for {
		select {
		case t := <-ResponseCH:
			if t == id {
				return
			}
		}
	}
}
