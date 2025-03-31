package commands

import (
	"crypto/ecdh"
	"encoding/base64"
)

// node send this to host server and host server will reply with 'setcipher' command;
// return (sharedKey []bytes, nodeName string, err)
func Connect(privateKey *ecdh.PrivateKey, id string, name string, peerPublicKeyBase64 string) ([]byte, string, error) {
	var peerPubKeyDecoded []byte = make([]byte, len(peerPublicKeyBase64))
	n, err := base64.RawStdEncoding.Decode(peerPubKeyDecoded, []byte(peerPublicKeyBase64))
	if err != nil {
		panic(err)
	}
	peerPubKeyDecoded = peerPubKeyDecoded[:n]
	peerPubKey, err := ecdh.P256().NewPublicKey(peerPubKeyDecoded)
	if err != nil {
		panic(err)
	}

	sharedKey, err := privateKey.ECDH(peerPubKey)
	if err != nil {
		panic(err)
	}

	return sharedKey, name, nil
}
