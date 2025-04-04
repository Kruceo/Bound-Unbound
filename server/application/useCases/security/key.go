package security

import (
	"crypto/ecdh"
	"crypto/rand"
	"log"
)

type GenKeysUseCase struct {
}

func (gk GenKeysUseCase) GenKeys() (*ecdh.PrivateKey, *ecdh.PublicKey) {
	priv, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	pub := priv.PublicKey()
	return priv, pub
}
