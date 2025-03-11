package main

import (
	"bufio"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"unbound-mngr-host/api/v1/handlers"
	"unbound-mngr-host/commands"
	"unbound-mngr-host/host"
	"unbound-mngr-host/memory"
	"unbound-mngr-host/security"
	"unbound-mngr-host/utils"

	"github.com/gorilla/websocket"
)

// Upgrader to upgrade HTTP to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections (Change this for security)
	},
}

func main() {
	isHost := len(os.Args) > 0 && os.Args[1] == "--host"
	if isHost {
		http.HandleFunc("/ws", handleWebSocket)

		http.HandleFunc("/v1/connections", handlers.ConnectionsHandler)

		http.HandleFunc("/v1/connections/{connection}/blocked", handlers.BlockAddressHandler)

		http.HandleFunc("/v1/connections/{connection}/redirects", handlers.RedirectAddressHandler)

		http.HandleFunc("/v1/connections/{connection}/reload", handlers.ReloadHandler)

		http.HandleFunc("/v1/connections/{connection}/confighash", handlers.ConfigHashHandler)

		fmt.Println("WebSocket server running on ws://localhost:8080/ws")

		// go func() {
		// 	scanner := bufio.NewScanner(os.Stdin)
		// 	for {
		// 		fmt.Print("\n>: ")
		// 		scanner.Scan()
		// 		text := scanner.Text()
		// 		HandleCommands(nil, "local "+text)
		// 	}
		// }()

		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Server error:", err)
		}
	} else {
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
}

func connectWebsocket() *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", host.MAIN_SERVER_ADDRESS), nil)
	if err != nil {
		fmt.Println("Connection error:", err)
		return nil
	}
	return conn
}

func HandleCommands(conn *websocket.Conn, str string, cipher **cipher.AEAD) {
	normalStr := str
	decrypted := false
	if cipher != nil && *cipher != nil {
		if strings.HasPrefix(str, "#$") {
			encodedStr, _ := strings.CutPrefix(normalStr, "#$")
			// fmt.Println("ENCODED", []byte(encodedStr))
			msg, err := security.DecipherMessageBase64Str(encodedStr, **cipher)
			if err != nil {
				fmt.Println(err)
				return
			}
			normalStr = string(msg)
			decrypted = true
		}
	}
	if decrypted {
		fmt.Println("[cmd decrypted]", normalStr)
	} else {
		fmt.Println("[cmd] ", str)
	}
	splt := strings.Split(normalStr, " ")
	if len(splt) < 2 {
		fmt.Println("wrong Syntax: (" + str + ")\nuse 'id' 'command'")
		return
	}
	id := splt[0]
	command := splt[1]
	args := splt[2:]

	if command == "connect" {
		err := commands.Connect(conn, id, args)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if command == "setcipher" && cipher != nil && *cipher == nil {
		c, err := commands.SetCipher(conn, id, args)
		if err != nil {
			fmt.Println(err)
			return
		}
		cptr := &c
		*cipher = cptr
		fmt.Println("new shared key added")
		memory.SetHost(conn, c)
	}
	if cipher != nil && decrypted {

		if command == "block" {
			err := commands.Block(conn, id, false, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if command == "unblock" {
			err := commands.Block(conn, id, true, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if command == "list" {
			err := commands.List(conn, id, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if command == "reload" {
			err := commands.ReloadConfig(conn, id)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if command == "redirect" {
			err := commands.AddRedirect(conn, id, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if command == "unredirect" {
			err := commands.RemoveRedirect(conn, id, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if command == "add" {
			err := commands.Add(conn, id, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected")

	for {
		// Read message from client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		splt := strings.Split(string(msg), " ")
		if len(splt) < 1 {
			fmt.Println("wrong syntax: " + string(msg))
			continue
		}
		cipher := memory.Connections[conn.RemoteAddr().String()].Cipher
		cptr := &cipher
		HandleCommands(conn, string(msg), &cptr)
		// err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}
