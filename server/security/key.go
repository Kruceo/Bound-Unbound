package security

import (
	"crypto/ecdh"
	"crypto/rand"
	"log"
)

var PrivateKey *ecdh.PrivateKey
var PublicKey *ecdh.PublicKey

func init() {
	priv, pub := GenKeys()
	PrivateKey = priv
	PublicKey = pub
}

func GenKeys() (*ecdh.PrivateKey, *ecdh.PublicKey) {
	priv, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	pub := priv.PublicKey()
	return priv, pub
}
