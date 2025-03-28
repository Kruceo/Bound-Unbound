package usecases

import (
	"crypto/cipher"
	"server2/application/entities"
	"server2/security"
	"strings"
)

type ParseCommandUseCase struct {
	Cipher *cipher.AEAD
}

func (d *ParseCommandUseCase) Execute(str string) (entities.Command, error) {
	message := str
	encrypted := false
	if *d.Cipher != nil {
		encodedStr, hasPrefix := strings.CutPrefix(str, "#$")
		if hasPrefix {
			msg, err := security.DecipherMessageBase64(encodedStr, *d.Cipher)
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

type SaveNodeUseCase struct {
	Repo entities.NodeRepository
}

func (r SaveNodeUseCase) Execute(node entities.Node) {
	r.Repo.Save(node)
}

// Get

type GetNodeUseCase struct {
	Repo entities.NodeRepository
}

func (r GetNodeUseCase) Execute(id string) *entities.Node {
	return r.Repo.Get(id)
}

// ids

type GetStoredNodesUseCase struct {
	Repo entities.NodeRepository
}

func (r GetStoredNodesUseCase) Execute() []string {
	return r.Repo.IDs()
}
