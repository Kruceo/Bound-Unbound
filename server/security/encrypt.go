package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"log"
)

func CreateCipher(secret []byte) cipher.AEAD {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		panic(err)
	}

	cipher, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	return cipher

	// encryptedMsg := mode.Seal(nil, nonce, []byte("vasco da gama"), nil)

	// fmt.Printf("%x\n", encryptedMsg)

	// decrypted, err := mode.Open(nil, nonce, encryptedMsg, nil)
	// if err != nil {
	// 	panic(err)
	// }
}

func RandomNonce() []byte {

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		log.Fatal(err)
	}
	return nonce
}
