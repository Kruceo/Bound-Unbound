//go:build !host
// +build !host

package main

import (
	"bufio"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"time"
	"unbound-mngr-host/host"
	"unbound-mngr-host/memory"
	"unbound-mngr-host/security"
	"unbound-mngr-host/utils"

	"github.com/gorilla/websocket"
)

func connectWebsocket() *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", host.MAIN_SERVER_ADDRESS), nil)
	if err != nil {
		fmt.Println("Connection error:", err)
		return nil
	}
	return conn
}

const IsHost = false

func Run() {
	name := utils.GetEnvOrDefault("NAME", fmt.Sprintf("%x", rand.Int()))
	fmt.Println(name)
	// Connect to WebSocket server
	var conn *websocket.Conn

	fmt.Println("trying connection")
	var cipher *cipher.AEAD
	go func() {
		for {
			if conn == nil {
				go func() {
					conn = connectWebsocket()
					if conn != nil {
						fmt.Println("sending and receiving keys")
						responseId := fmt.Sprintf("%x", rand.Int())
						var encodedPublicKey = base64.RawStdEncoding.EncodeToString(security.PublicKey.Bytes())
						conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s connect %s %s", responseId, encodedPublicKey, name)))

						fmt.Println("waiting for response... (" + responseId + ")")
						memory.WaitForResponse(responseId)

					}
				}()
				time.Sleep(1 * time.Second)
				continue
			}

			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Read error:", err)
				conn.Close()
				conn = nil
				cipher = nil
				memory.RemoveHost()
			} else {
				HandleCommands(conn, string(msg), &cipher)
			}

		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\n>: ")
		scanner.Scan()
		text := scanner.Text()
		HandleCommands(conn, "local "+text, &cipher)
	}
}
