package usecases

import (
	"crypto/cipher"
	"fmt"
	"server2/application/entities"
	"server2/application/infrastructure"

	"github.com/gorilla/websocket"
)

type NodePersistenceUseCase struct {
	repo infrastructure.NodeRepository
}

func NewNodePersistenceUseCase(repo infrastructure.NodeRepository) *NodePersistenceUseCase {
	return &NodePersistenceUseCase{repo: repo}
}

func (r *NodePersistenceUseCase) Save(id, name string, conn *websocket.Conn, cipher *cipher.AEAD) (string, error) {
	return r.repo.Save(id, name, conn, cipher)
}

// Get

func (r *NodePersistenceUseCase) Get(id string) *entities.Node {
	return r.repo.Get(id)
}

// ids

func (r NodePersistenceUseCase) IDs() []string {
	return r.repo.IDs()
}

// delete

func (r *NodePersistenceUseCase) Delete(id string) error {
	return r.repo.Delete(id)
}

func (uc *NodePersistenceUseCase) GetOrCreate(nodeID string, conn *websocket.Conn) (*entities.Node, error) {
	node := uc.repo.Get(nodeID)
	if node == nil {

		_, err := uc.repo.Save(nodeID, node.Name, conn, nil)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return node, nil
}

func (r *NodePersistenceUseCase) FindOneByRemoteAddress(remoteAddr string) (*entities.Node, error) {
	return r.repo.FindOneByRemoteAddress(remoteAddr)
}
