package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
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
}

func RandomNonce() []byte {

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		log.Fatal(err)
	}
	return nonce
}

func DecipherMessageBase64Str(str string, cipher cipher.AEAD) ([]byte, error) {
	var decodedStr []byte = make([]byte, len(str)+12)
	n, err := base64.RawStdEncoding.Decode(decodedStr, []byte(str))
	if err != nil {
		return nil, err
	}
	decodedStr = decodedStr[:n]
	content := decodedStr[0 : n-12]
	nonce := decodedStr[n-12:]
	result, err := cipher.Open(nil, nonce, content, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}
