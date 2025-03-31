package usecases

import (
	"crypto/cipher"
	"fmt"
	"server2/application/entities"
	"server2/security"
	"strings"

	"github.com/gorilla/websocket"
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

func (r SaveNodeUseCase) Execute(Conn *websocket.Conn, Name string, Cipher cipher.AEAD) (string, error) {
	return r.Repo.Save(entities.Node{Conn: Conn, Name: Name, Cipher: Cipher})
}

type CreateNodeUseCase struct{}

func (r CreateNodeUseCase) Execute(conn *websocket.Conn, name string, cipher cipher.AEAD) (*entities.Node, error) {
	node := entities.Node{Conn: conn, Name: name, Cipher: cipher}
	if len(name) == 0 {
		return nil, fmt.Errorf("bad name: %s", name)
	}
	if conn == nil {
		return nil, fmt.Errorf("bad connection: %v", conn)
	}
	return &node, nil
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

// delete

type DeleteNodeUseCase struct {
	Repo entities.NodeRepository
}

func (r DeleteNodeUseCase) Execute(id string) error {
	return r.Repo.Delete(id)
}

type GetOrCreateUseCase struct {
	Repo entities.NodeRepository
}

func (uc *GetOrCreateUseCase) Execute(nodeID string, conn *websocket.Conn) (*entities.Node, error) {
	node := uc.Repo.Get(nodeID)
	if node == nil {
		node = &entities.Node{Conn: conn, Name: "nameless", Cipher: nil}
		_, err := uc.Repo.Save(*node)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return node, nil
}

type CipherMessageUseCase struct {
}

func (c *CipherMessageUseCase) Execute(message string, cipher *cipher.AEAD) ([]byte, error) {
	base64Msg := security.CipherMessageBase64(message, *cipher)
	return append([]byte("#$"), base64Msg...), nil
}

func NewCipherMessageUseCase(cipher cipher.AEAD) CipherMessageUseCase {
	return CipherMessageUseCase{}
}
