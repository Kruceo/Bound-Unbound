package commands

import (
	"crypto/cipher"
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
	"server2/security"
	"strings"
)

// node send this to host server and host server will reply with 'setcipher' command;
// return (sharedKey []bytes, nodeName string, err)
func Connect(id string, args []string) ([]byte, string, error) {
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

	// go func() {
	// 	ticker := time.NewTicker(1 * time.Second)
	// 	defer ticker.Stop()
	// 	for range ticker.C {
	// 		err := conn.WriteMessage(websocket.PingMessage, nil)
	// 		if err != nil {
	// 			delete(memory.Connections, conn.RemoteAddr().String())
	// 			fmt.Println(name, "disconnected")
	// 			break
	// 		}
	// 	}
	// }()
	// encodedPublicKey := base64.RawStdEncoding.EncodeToString(security.PublicKey.Bytes())
	// conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("_ setcipher %s", encodedPublicKey)))
	return sharedKey, name, nil
}

// return aes cipher, error
func SetCipher(id string, args []string) (cipher.AEAD, error) {
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
