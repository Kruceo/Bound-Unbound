//go:build host
// +build host

package commands

import (
	"crypto/cipher"
	"fmt"
	"strings"
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

	if command == "connect" {
		err := Connect(conn, id, args)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if cipher != nil && decrypted {
		if command == "add" {
			err := Add(conn, id, args)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
