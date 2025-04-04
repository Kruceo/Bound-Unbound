package security_test

// func TestEncrypt(t *testing.T) {
// 	original := "real message; terra; logging; redirects; #@!¨&*()"
// 	secret := make([]byte, 22)
// 	_, err := rand.Read(secret)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	fmt.Println("secret =", secret)
// 	cipher := security.CreateCipher(secret)
// 	nonce := security.RandomNonce()
// 	fmt.Println("nonce =", nonce)
// 	encrypted := cipher.Seal(nil, nonce, []byte(original), nil)
// 	result, err := cipher.Open(nil, nonce, encrypted, nil)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if string(result) != original {
// 		t.Errorf("result is not equal original: [%s] VS [%s]", result, original)
// 	}
// }

// func TestCipherMessageBase64(t *testing.T) {
// 	original := "Mr. Gordon is landing in Mars."
// 	cip := security.CreateCipher([]byte("logaritm"))

// 	encodedMessage := security.CipherMessageBase64(original, cip)
// 	decodedMessage, err := security.DecipherMessageBase64(string(encodedMessage), cip)

// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	fmt.Printf("original = %s\n", original)
// 	fmt.Printf("encoded = %x\n", encodedMessage)
// 	fmt.Printf("result = %s\n", decodedMessage)

// 	if original != string(decodedMessage) {
// 		t.Errorf("result is not equal original: [%s] VS [%s]", decodedMessage, original)
// 	}
// }
