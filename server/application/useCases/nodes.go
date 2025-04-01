package usecases

import (
	"crypto/cipher"
	"fmt"
	"server2/application/entities"
	"server2/application/infrastructure"

	"github.com/gorilla/websocket"
)

type SaveNodeUseCase struct {
	Repo *infrastructure.NodeRepository
}

func (r *SaveNodeUseCase) Execute(Conn *websocket.Conn, Name string, Cipher cipher.AEAD) (string, error) {
	return (*r.Repo).Save(entities.Node{Conn: Conn, Name: Name, Cipher: Cipher})
}

type CreateNodeUseCase struct{}

func (r *CreateNodeUseCase) Execute(conn *websocket.Conn, name string, cipher cipher.AEAD) (*entities.Node, error) {
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
	Repo *infrastructure.NodeRepository
}

func (r *GetNodeUseCase) Execute(id string) *entities.Node {
	return (*r.Repo).Get(id)
}

// ids

type GetStoredNodesUseCase struct {
	Repo *infrastructure.NodeRepository
}

func (r GetStoredNodesUseCase) Execute() []string {
	return (*r.Repo).IDs()
}

// delete

type DeleteNodeUseCase struct {
	Repo *infrastructure.NodeRepository
}

func (r *DeleteNodeUseCase) Execute(id string) error {
	return (*r.Repo).Delete(id)
}

type GetOrCreateUseCase struct {
	Repo *infrastructure.NodeRepository
}

func (uc *GetOrCreateUseCase) Execute(nodeID string, conn *websocket.Conn) (*entities.Node, error) {
	node := (*uc.Repo).Get(nodeID)
	if node == nil {
		node = &entities.Node{Conn: conn, Name: "nameless", Cipher: nil}
		_, err := (*uc.Repo).Save(*node)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return node, nil
}
