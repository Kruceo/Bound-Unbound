package main

import (
	"crypto/cipher"
	"fmt"
	"strings"
	v1 "unbound-mngr-host/api/v1"
	"unbound-mngr-host/commands"
	"unbound-mngr-host/host"
	"unbound-mngr-host/memory"
	"unbound-mngr-host/security"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func HandleCommands(conn *websocket.Conn, str string, cipher **cipher.AEAD) {
	normalStr := str
	decrypted := false
	if cipher != nil && *cipher != nil {
		encodedStr, hasPrefix := strings.CutPrefix(normalStr, "#$")
		if hasPrefix {
			msg, err := security.DecipherMessageBase64(encodedStr, **cipher)
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

	if IsHost { // HOST
		if command == "connect" {
			err := commands.Connect(conn, id, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if cipher != nil && decrypted {
			if command == "add" {
				err := commands.Add(conn, id, args)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	} else { // NODE
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

}

func init() {
	godotenv.Load(".env")
	host.InitLocals()
	v1.InitAuth()
}

func main() {
	Run()
}
