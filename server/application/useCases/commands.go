package usecases

import (
	"crypto/cipher"
	"server2/application/entities"
	"server2/application/useCases/security"
	"strings"
)

type ParseCommandUseCase struct {
	Cipher         *cipher.AEAD
	ciphersUseCase security.CiphersUseCase
}

func (d *ParseCommandUseCase) Execute(str string) (entities.Command, error) {
	message := str
	encrypted := false
	if d.Cipher != nil {
		encodedStr, hasPrefix := strings.CutPrefix(str, "#$")
		if hasPrefix {
			msg, err := d.ciphersUseCase.DecipherMessageBase64(encodedStr, *d.Cipher)
			if err != nil {
				return entities.Command{}, err
			}
			message = string(msg)
			encrypted = true
		}
	}

	spl := strings.Split(message, " ")

	id, entry, args := spl[0], spl[1], spl[2:]

	newArgs := []string{}
	mode := false

	for _, v := range args {
		if strings.HasPrefix(v, "\"") {
			newArgs = append(newArgs, "")
			mode = true
		} else if strings.HasSuffix(v, "\"") && !strings.HasSuffix(v, "\\\"") {
			newArgs[len(newArgs)-1] += v
			mode = false
			continue
		}

		if mode && len(newArgs) > 0 {
			newArgs[len(newArgs)-1] += v + " "
		} else {
			newArgs = append(newArgs, v)
		}

	}
	return entities.Command{Entry: entry, Args: newArgs, Id: id, IsEncrypted: encrypted, Raw: str}, nil
}

// nodes
// Save

type CipherCommandMessageUseCase struct {
	ciphersUseCase security.CiphersUseCase
}

func (c *CipherCommandMessageUseCase) Execute(message string, cipher *cipher.AEAD) ([]byte, error) {
	base64Msg := c.ciphersUseCase.CipherMessageBase64(message, *cipher)
	return append([]byte("#$"), base64Msg...), nil
}

func NewCipherMessageUseCase() CipherCommandMessageUseCase {
	return CipherCommandMessageUseCase{}
}
