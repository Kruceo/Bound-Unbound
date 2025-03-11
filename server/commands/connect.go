package commands

import (
	"crypto/cipher"
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
	"unbound-mngr-host/memory"
	"unbound-mngr-host/security"

	"github.com/gorilla/websocket"
)

func Connect(conn *websocket.Conn, id string, args []string) error {
	// fmt.Println(args)
	name := strings.Join(args[1:], " ")

	// add connection(arg0) <id>(arg1)
	var peerPubKeyDecoded []byte = make([]byte, len(args[0]))
	n, err := base64.RawStdEncoding.Decode(peerPubKeyDecoded, []byte(args[0]))
	if err != nil {
		panic(err)
	}
	peerPubKeyDecoded = peerPubKeyDecoded[:n]
	peerPubKey, err := ecdh.P256().NewPublicKey(peerPubKeyDecoded)
	if err != nil {
		panic(err)
	}

	sharedKey, err := security.PrivateKey.ECDH(peerPubKey)
	if err != nil {
		panic(err)
	}
	memory.Connections[conn.RemoteAddr().String()] = memory.Client{Conn: conn, Name: name, Cipher: security.CreateCipher(sharedKey)}
	fmt.Println(name, conn.LocalAddr().String(), "connected")

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				delete(memory.Connections, conn.RemoteAddr().String())
				fmt.Println(name, "disconnected")
				break
			}
		}
	}()
	encodedPublicKey := base64.RawStdEncoding.EncodeToString(security.PublicKey.Bytes())
	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("_ setcipher %s", encodedPublicKey)))
	return nil
}

func SetCipher(conn *websocket.Conn, id string, args []string) (cipher.AEAD, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("wrong syntax: (%v)\nuse setcipher", args)
	}
	hostPublicKeyEncoded := args[0]

	var peerPubKeyDecoded []byte = make([]byte, len(hostPublicKeyEncoded))
	n, err := base64.RawStdEncoding.Decode(peerPubKeyDecoded, []byte(hostPublicKeyEncoded))
	if err != nil {
		panic(err)
	}
	peerPubKeyDecoded = peerPubKeyDecoded[:n]
	peerPubKey, err := ecdh.P256().NewPublicKey(peerPubKeyDecoded)
	if err != nil {
		panic(err)
	}
	sharedKey, err := security.PrivateKey.ECDH(peerPubKey)
	if err != nil {
		panic(err)
	}
	cipher := security.CreateCipher(sharedKey)
	return cipher, nil
}
