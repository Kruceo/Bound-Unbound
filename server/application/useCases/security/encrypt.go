package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

type CiphersUseCase struct {
}

func (c CiphersUseCase) CreateCipher(secret []byte) cipher.AEAD {
	if lenght := len(secret); lenght < 24 {
		toAdd := make([]byte, 24-lenght)
		secret = append(secret, toAdd...)
	} else if lenght := len(secret); lenght > 24 && lenght <= 32 {
		toAdd := make([]byte, 32-lenght)
		secret = append(secret, toAdd...)
	} else {
		panic(fmt.Errorf("secret is greater than 32 bytes: %v", secret))
	}

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

func (c CiphersUseCase) RandomNonce() []byte {

	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		log.Fatal(err)
	}
	return nonce
}

func (c CiphersUseCase) DecipherMessageBase64(str string, cipher cipher.AEAD) ([]byte, error) {
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

func (c CiphersUseCase) CipherMessageBase64(str string, cipher cipher.AEAD) []byte {
	nonce := c.RandomNonce()
	encryptedContent := cipher.Seal(nil, nonce, []byte(str), nil)
	formatedMsg := encryptedContent
	formatedMsg = append(formatedMsg, nonce...)
	base64Msg := base64.RawStdEncoding.EncodeToString(formatedMsg)
	return []byte(base64Msg)
}
