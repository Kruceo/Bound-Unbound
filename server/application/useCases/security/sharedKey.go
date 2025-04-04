package security

import (
	"crypto/ecdh"
	"encoding/base64"
)

type CreateSharedKeyUseCase struct {
	privateKey *ecdh.PrivateKey
}

// node send this to host server and host server will reply with other 'connect' command;
// return (sharedKey []bytes, nodeName string, err)
func (uc *CreateSharedKeyUseCase) Execute(peerPublicKeyBase64 string) ([]byte, error) {
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

	sharedKey, err := uc.privateKey.ECDH(peerPubKey)
	if err != nil {
		panic(err)
	}

	return sharedKey, nil
}

func NewCreateSharedKeyUseCase(privKey ecdh.PrivateKey) CreateSharedKeyUseCase {
	return CreateSharedKeyUseCase{privateKey: &privKey}
}
