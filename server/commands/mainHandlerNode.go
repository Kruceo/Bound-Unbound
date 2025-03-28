//go:build !host
// +build !host

package commands

import (
	"crypto/cipher"
	"fmt"
	"strings"
	"unbound-mngr-host/memory"
	"unbound-mngr-host/security"

	"github.com/gorilla/websocket"
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

	if command == "setcipher" && cipher != nil && *cipher == nil {
		c, err := SetCipher(conn, id, args)
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
			err := Block(conn, id, false, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if command == "unblock" {
			err := Block(conn, id, true, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if command == "list" {
			err := List(conn, id, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if command == "reload" {
			err := ReloadConfig(conn, id)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if command == "redirect" {
			err := AddRedirect(conn, id, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if command == "unredirect" {
			err := RemoveRedirect(conn, id, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if command == "add" {
			err := Add(conn, id, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
